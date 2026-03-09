import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:mocktail/mocktail.dart';
import 'package:dio/dio.dart';

import 'package:open_talk/core/theme/app_theme.dart';
import 'package:open_talk/features/auth/data/repositories/auth_repository.dart';
import 'package:open_talk/features/auth/domain/providers/auth_provider.dart';
import 'package:open_talk/features/auth/presentation/screens/login_screen.dart';
import 'package:open_talk/shared/services/auth_storage.dart';
import 'package:open_talk/shared/widgets/app_button.dart';
import 'package:open_talk/shared/widgets/app_text_field.dart';

// ---------------------------------------------------------------------------
// Proper mocks using mocktail
// ---------------------------------------------------------------------------
class MockAuthRepository extends Mock implements AuthRepository {}
class MockAuthStorage extends Mock implements AuthStorage {}
class MockDio extends Mock implements Dio {}

// ---------------------------------------------------------------------------
// Stub AuthNotifier that accepts a pre-set AuthState and never touches
// real network or storage.
// ---------------------------------------------------------------------------
class StubAuthNotifier extends AuthNotifier {
  StubAuthNotifier(AuthState initialState, {
    required AuthRepository authRepository,
    required AuthStorage authStorage,
  }) : super(
          authRepository: authRepository,
          authStorage: authStorage,
        ) {
    state = initialState;
  }
}

// ---------------------------------------------------------------------------
// Helper
// ---------------------------------------------------------------------------
Widget _buildLoginScreen({AuthState? authState}) {
  final mockRepo = MockAuthRepository();
  final mockStorage = MockAuthStorage();

  // Stub the hasToken call so checkAuthStatus (if triggered) doesn't throw
  when(() => mockStorage.hasToken()).thenAnswer((_) async => false);
  when(() => mockStorage.getAccessToken()).thenAnswer((_) async => null);

  final state = authState ?? const AuthState(status: AuthStatus.unauthenticated);

  return ProviderScope(
    overrides: [
      authStateProvider.overrideWith(
        (ref) => StubAuthNotifier(
          state,
          authRepository: mockRepo,
          authStorage: mockStorage,
        ),
      ),
    ],
    child: MaterialApp(
      theme: AppTheme.light,
      home: const LoginScreen(),
    ),
  );
}

void main() {
  group('LoginScreen', () {
    group('renders correctly', () {
      testWidgets('shows email text field with correct label', (tester) async {
        await tester.pumpWidget(_buildLoginScreen());
        await tester.pump();

        expect(find.text('이메일'), findsOneWidget);
      });

      testWidgets('shows password text field with correct label', (tester) async {
        await tester.pumpWidget(_buildLoginScreen());
        await tester.pump();

        expect(find.text('비밀번호'), findsOneWidget);
      });

      testWidgets('shows login button', (tester) async {
        await tester.pumpWidget(_buildLoginScreen());
        await tester.pump();

        // '로그인' appears as button label and possibly in other text
        expect(find.text('로그인'), findsWidgets);
      });

      testWidgets('shows register link text button', (tester) async {
        await tester.pumpWidget(_buildLoginScreen());
        await tester.pump();

        expect(find.text('회원가입'), findsOneWidget);
      });

      testWidgets('shows app name', (tester) async {
        await tester.pumpWidget(_buildLoginScreen());
        await tester.pump();

        expect(find.text('OpenTalk'), findsOneWidget);
      });

      testWidgets('renders two AppTextField widgets', (tester) async {
        await tester.pumpWidget(_buildLoginScreen());
        await tester.pump();

        expect(find.byType(AppTextField), findsNWidgets(2));
      });

      testWidgets('renders AppButton for login action', (tester) async {
        await tester.pumpWidget(_buildLoginScreen());
        await tester.pump();

        expect(find.byType(AppButton), findsOneWidget);
      });
    });

    group('validation on empty submit', () {
      testWidgets('shows email validation error when email field is empty', (tester) async {
        await tester.pumpWidget(_buildLoginScreen());
        await tester.pump();

        await tester.tap(find.byType(ElevatedButton));
        await tester.pump();

        expect(find.text('이메일을 입력해주세요'), findsOneWidget);
      });

      testWidgets('shows password validation error when password field is empty', (tester) async {
        await tester.pumpWidget(_buildLoginScreen());
        await tester.pump();

        // Enter valid email but leave password empty
        await tester.enterText(
          find.byType(TextFormField).first,
          'alice@example.com',
        );
        await tester.tap(find.byType(ElevatedButton));
        await tester.pump();

        expect(find.text('비밀번호를 입력해주세요'), findsOneWidget);
      });

      testWidgets('shows email format error for invalid email', (tester) async {
        await tester.pumpWidget(_buildLoginScreen());
        await tester.pump();

        await tester.enterText(find.byType(TextFormField).first, 'notanemail');
        await tester.tap(find.byType(ElevatedButton));
        await tester.pump();

        expect(find.text('올바른 이메일 형식이 아닙니다'), findsOneWidget);
      });
    });

    group('loading state', () {
      testWidgets('button shows CircularProgressIndicator when status is loading', (tester) async {
        await tester.pumpWidget(
          _buildLoginScreen(
            authState: const AuthState(status: AuthStatus.loading),
          ),
        );
        await tester.pump();

        expect(find.byType(CircularProgressIndicator), findsOneWidget);
      });
    });
  });
}
