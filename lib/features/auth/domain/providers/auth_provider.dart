import 'package:dio/dio.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../../core/network/api_client.dart';
import '../../../../shared/services/auth_storage.dart';
import '../../data/models/user_model.dart';
import '../../data/repositories/auth_repository.dart';

// Repository
final authRepositoryProvider = Provider<AuthRepository>((ref) {
  return AuthRepository(
    dio: ref.read(apiClientProvider),
    authStorage: ref.read(authStorageProvider),
  );
});

// Auth State
enum AuthStatus { initial, authenticated, unauthenticated, loading }

class AuthState {
  final AuthStatus status;
  final UserModel? user;
  final String? error;

  const AuthState({
    this.status = AuthStatus.initial,
    this.user,
    this.error,
  });

  AuthState copyWith({
    AuthStatus? status,
    UserModel? user,
    String? error,
  }) {
    return AuthState(
      status: status ?? this.status,
      user: user ?? this.user,
      error: error,
    );
  }
}

class AuthNotifier extends StateNotifier<AuthState> {
  final AuthRepository _authRepository;
  final AuthStorage _authStorage;

  AuthNotifier({
    required AuthRepository authRepository,
    required AuthStorage authStorage,
  })  : _authRepository = authRepository,
        _authStorage = authStorage,
        super(const AuthState());

  Future<void> checkAuthStatus() async {
    state = state.copyWith(status: AuthStatus.loading);

    final hasToken = await _authStorage.hasToken();
    if (!hasToken) {
      state = state.copyWith(status: AuthStatus.unauthenticated);
      return;
    }

    try {
      final user = await _authRepository.getMe();
      state = state.copyWith(
        status: AuthStatus.authenticated,
        user: user,
      );
    } on DioException {
      state = state.copyWith(status: AuthStatus.unauthenticated);
      await _authStorage.clearAll();
    }
  }

  Future<void> login({
    required String email,
    required String password,
  }) async {
    state = state.copyWith(status: AuthStatus.loading, error: null);

    try {
      final authResponse = await _authRepository.login(
        email: email,
        password: password,
      );
      state = state.copyWith(
        status: AuthStatus.authenticated,
        user: authResponse.user,
      );
    } on DioException catch (e) {
      final message = _extractErrorMessage(e);
      state = state.copyWith(
        status: AuthStatus.unauthenticated,
        error: message,
      );
    }
  }

  Future<void> register({
    required String email,
    required String password,
    required String nickname,
  }) async {
    state = state.copyWith(status: AuthStatus.loading, error: null);

    try {
      final authResponse = await _authRepository.register(
        email: email,
        password: password,
        nickname: nickname,
      );
      state = state.copyWith(
        status: AuthStatus.authenticated,
        user: authResponse.user,
      );
    } on DioException catch (e) {
      final message = _extractErrorMessage(e);
      state = state.copyWith(
        status: AuthStatus.unauthenticated,
        error: message,
      );
    }
  }

  Future<void> logout() async {
    await _authRepository.logout();
    state = const AuthState(status: AuthStatus.unauthenticated);
  }

  void clearError() {
    state = state.copyWith(error: null);
  }

  String _extractErrorMessage(DioException e) {
    if (e.response?.data is Map<String, dynamic>) {
      final data = e.response!.data as Map<String, dynamic>;
      return data['message'] as String? ?? '오류가 발생했습니다';
    }
    if (e.type == DioExceptionType.connectionTimeout ||
        e.type == DioExceptionType.receiveTimeout) {
      return '서버 연결 시간이 초과되었습니다';
    }
    if (e.type == DioExceptionType.connectionError) {
      return '서버에 연결할 수 없습니다';
    }
    return '오류가 발생했습니다';
  }
}

final authStateProvider =
    StateNotifierProvider<AuthNotifier, AuthState>((ref) {
  return AuthNotifier(
    authRepository: ref.read(authRepositoryProvider),
    authStorage: ref.read(authStorageProvider),
  );
});

final currentUserProvider = Provider<UserModel?>((ref) {
  return ref.watch(authStateProvider).user;
});
