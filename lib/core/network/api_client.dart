import 'package:dio/dio.dart';
import 'package:flutter/foundation.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../constants/app_constants.dart';
import '../../shared/services/auth_storage.dart';

final apiClientProvider = Provider<Dio>((ref) {
  final dio = Dio(
    BaseOptions(
      baseUrl: AppConstants.apiBaseUrl,
      connectTimeout: const Duration(seconds: 10),
      receiveTimeout: const Duration(seconds: 10),
      headers: {
        'Content-Type': 'application/json',
        'Accept': 'application/json',
      },
    ),
  );

  final authStorage = ref.read(authStorageProvider);

  dio.interceptors.add(AuthInterceptor(authStorage));
  dio.interceptors.add(ErrorInterceptor(ref));

  if (kDebugMode) {
    dio.interceptors.add(
      LogInterceptor(
        requestBody: true,
        responseBody: true,
        logPrint: (o) => debugPrint(o.toString()),
      ),
    );
  }

  return dio;
});

class AuthInterceptor extends Interceptor {
  final AuthStorage _authStorage;

  AuthInterceptor(this._authStorage);

  @override
  Future<void> onRequest(
    RequestOptions options,
    RequestInterceptorHandler handler,
  ) async {
    final token = await _authStorage.getAccessToken();
    if (token != null) {
      options.headers['Authorization'] = 'Bearer $token';
    }
    handler.next(options);
  }
}

class ErrorInterceptor extends Interceptor {
  final Ref _ref;
  bool _isRefreshing = false;

  ErrorInterceptor(this._ref);

  @override
  void onError(DioException err, ErrorInterceptorHandler handler) async {
    if (err.response?.statusCode == 401 && !_isRefreshing) {
      _isRefreshing = true;
      try {
        final authStorage = _ref.read(authStorageProvider);
        final refreshToken = await authStorage.getRefreshToken();
        if (refreshToken != null) {
          final dio = Dio(BaseOptions(
            baseUrl: AppConstants.apiBaseUrl,
            headers: {'Content-Type': 'application/json'},
          ));
          final response = await dio.post(
            '/auth/refresh',
            data: {'token': refreshToken},
          );
          final newToken =
              (response.data as Map<String, dynamic>)['data']['token'] as String;
          await authStorage.saveAccessToken(newToken);
          _isRefreshing = false;

          // Retry original request
          final opts = err.requestOptions;
          opts.headers['Authorization'] = 'Bearer $newToken';
          final retryResponse = await dio.fetch(opts);
          return handler.resolve(retryResponse);
        }
      } catch (_) {
        // Refresh failed, clear auth
        await _ref.read(authStorageProvider).clearAll();
      }
      _isRefreshing = false;
    }
    handler.next(err);
  }
}
