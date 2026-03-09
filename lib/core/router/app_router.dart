import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../features/auth/domain/providers/auth_provider.dart';
import '../../features/auth/presentation/screens/login_screen.dart';
import '../../features/auth/presentation/screens/register_screen.dart';
import '../../features/auth/presentation/screens/splash_screen.dart';
import '../../features/chat/presentation/screens/chat_room_screen.dart';
import '../../features/home/presentation/screens/add_friend_screen.dart';
import '../../features/home/presentation/screens/friend_requests_screen.dart';
import '../../features/home/presentation/screens/home_screen.dart';
import '../../features/home/presentation/screens/password_change_screen.dart';
import '../../features/home/presentation/screens/profile_edit_screen.dart';
import '../../features/home/presentation/screens/profile_screen.dart';

class AuthNotifierListenable extends ChangeNotifier {
  AuthNotifierListenable(Ref ref) {
    ref.listen<AuthState>(authStateProvider, (prev, next) {
      notifyListeners();
    });
  }
}

final _authListenableProvider = Provider<AuthNotifierListenable>((ref) {
  return AuthNotifierListenable(ref);
});

final routerProvider = Provider<GoRouter>((ref) {
  final authListenable = ref.read(_authListenableProvider);

  return GoRouter(
    initialLocation: '/splash',
    debugLogDiagnostics: true,
    refreshListenable: authListenable,
    redirect: (context, state) {
      final authState = ref.read(authStateProvider);
      final isAuthenticated = authState.status == AuthStatus.authenticated;
      final isLoading = authState.status == AuthStatus.loading ||
          authState.status == AuthStatus.initial;
      final currentPath = state.matchedLocation;

      // Allow splash screen during loading
      if (isLoading && currentPath == '/splash') return null;

      // Auth pages
      final isAuthPage = currentPath == '/login' ||
          currentPath == '/register' ||
          currentPath == '/splash';

      if (!isAuthenticated && !isAuthPage) return '/login';
      if (isAuthenticated && isAuthPage) return '/home';

      return null;
    },
    routes: [
      GoRoute(
        path: '/splash',
        builder: (context, state) => const SplashScreen(),
      ),
      GoRoute(
        path: '/login',
        builder: (context, state) => const LoginScreen(),
      ),
      GoRoute(
        path: '/register',
        builder: (context, state) => const RegisterScreen(),
      ),
      GoRoute(
        path: '/home',
        builder: (context, state) => const HomeScreen(),
      ),
      GoRoute(
        path: '/chat/:roomId',
        builder: (context, state) {
          final roomId = state.pathParameters['roomId']!;
          return ChatRoomScreen(roomId: roomId);
        },
      ),
      GoRoute(
        path: '/profile',
        builder: (context, state) => const ProfileScreen(),
      ),
      GoRoute(
        path: '/profile/edit',
        builder: (context, state) => const ProfileEditScreen(),
      ),
      GoRoute(
        path: '/password/change',
        builder: (context, state) => const PasswordChangeScreen(),
      ),
      GoRoute(
        path: '/friends/add',
        builder: (context, state) => const AddFriendScreen(),
      ),
      GoRoute(
        path: '/friends/requests',
        builder: (context, state) => const FriendRequestsScreen(),
      ),
    ],
    errorBuilder: (context, state) => Scaffold(
      body: Center(
        child: Text('페이지를 찾을 수 없습니다: ${state.error}'),
      ),
    ),
  );
});
