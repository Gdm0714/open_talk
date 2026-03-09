import 'package:dio/dio.dart';

import '../../../../core/network/api_endpoints.dart';
import '../../../../shared/services/auth_storage.dart';
import '../models/auth_response_model.dart';
import '../models/user_model.dart';

class AuthRepository {
  final Dio _dio;
  final AuthStorage _authStorage;

  AuthRepository({
    required Dio dio,
    required AuthStorage authStorage,
  })  : _dio = dio,
        _authStorage = authStorage;

  Future<AuthResponseModel> login({
    required String email,
    required String password,
  }) async {
    final response = await _dio.post(
      ApiEndpoints.login,
      data: {
        'email': email,
        'password': password,
      },
    );

    final authResponse = AuthResponseModel.fromJson(
      (response.data as Map<String, dynamic>)['data'] as Map<String, dynamic>,
    );

    await _saveTokens(authResponse);
    return authResponse;
  }

  Future<AuthResponseModel> register({
    required String email,
    required String password,
    required String nickname,
  }) async {
    final response = await _dio.post(
      ApiEndpoints.register,
      data: {
        'email': email,
        'password': password,
        'nickname': nickname,
      },
    );

    final authResponse = AuthResponseModel.fromJson(
      (response.data as Map<String, dynamic>)['data'] as Map<String, dynamic>,
    );

    await _saveTokens(authResponse);
    return authResponse;
  }

  Future<UserModel> getMe() async {
    final response = await _dio.get(ApiEndpoints.me);
    return UserModel.fromJson(
      (response.data as Map<String, dynamic>)['data'] as Map<String, dynamic>,
    );
  }

  Future<void> refreshToken() async {
    final refreshToken = await _authStorage.getRefreshToken();
    if (refreshToken == null) {
      throw Exception('No refresh token');
    }

    final response = await _dio.post(
      ApiEndpoints.refreshToken,
      data: {'refresh_token': refreshToken},
    );

    final data =
        (response.data as Map<String, dynamic>)['data'] as Map<String, dynamic>;
    await _authStorage.saveAccessToken(data['token'] as String);
  }

  Future<void> logout() async {
    try {
      await _dio.post(ApiEndpoints.logout);
    } catch (_) {
      // Ignore error on logout
    } finally {
      await _authStorage.clearAll();
    }
  }

  Future<void> _saveTokens(AuthResponseModel authResponse) async {
    await _authStorage.saveAccessToken(authResponse.accessToken);
    await _authStorage.saveRefreshToken(authResponse.refreshToken);
    await _authStorage.saveUserId(authResponse.user.id);
  }
}
