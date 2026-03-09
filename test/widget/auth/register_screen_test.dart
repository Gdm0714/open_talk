import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:mocktail/mocktail.dart';

import 'package:open_talk/core/theme/app_theme.dart';
import 'package:open_talk/features/auth/data/repositories/auth_repository.dart';
import 'package:open_talk/features/auth/domain/providers/auth_provider.dart';
import 'package:open_talk/features/auth/presentation/screens/register_screen.dart';
import 'package:open_talk/shared/services/auth_storage.dart';
import 'package:open_talk/shared/widgets/app_button.dart';
import 'package:open_talk/shared/widgets/app_text_field.dart';

// ---------------------------------------------------------------------------
// Proper mocks using mocktail
// ---------------------------------------------------------------------------
class MockAuthRepository extends Mock implements AuthRepository {}
class MockAuthStorage extends Mock implements AuthStorage {}

// ---------------------------------------------------------------------------
// Stub notifier that holds a pre-set AuthState
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
Widget _buildRegisterScreen({AuthState? authState}) {
  final mockRepo = MockAuthRepository();
  final mockStorage = MockAuthStorage();

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
      home: const RegisterScreen(),
    ),
  );
}

void main() {
  group('RegisterScreen', () {
    group('renders correctly', () {
      testWidgets('shows nickname text field', (tester) async {
        await tester.pumpWidget(_buildRegisterScreen());
        await tester.pump();

        expect(find.text('닉네임'), findsOneWidget);
      });

      testWidgets('shows email text field', (tester) async {
        await tester.pumpWidget(_buildRegisterScreen());
        await tester.pump();

        expect(find.text('이메일'), findsOneWidget);
      });

      testWidgets('shows password text field', (tester) async {
        await tester.pumpWidget(_buildRegisterScreen());
        await tester.pump();

        expect(find.text('비밀번호'), findsOneWidget);
      });

      testWidgets('renders exactly three AppTextField widgets', (tester) async {
        await tester.pumpWidget(_buildRegisterScreen());
        await tester.pump();

        expect(find.byType(AppTextField), findsNWidgets(3));
      });

      testWidgets('shows register button', (tester) async {
        await tester.pumpWidget(_buildRegisterScreen());
        await tester.pump();

        expect(find.text('회원가입'), findsWidgets);
      });

      testWidgets('shows login link for existing users', (tester) async {
        await tester.pumpWidget(_buildRegisterScreen());
        await tester.pump();

        expect(find.text('로그인'), findsOneWidget);
      });

      testWidgets('renders AppButton for register action', (tester) async {
        await tester.pumpWidget(_buildRegisterScreen());
        await tester.pump();

        expect(find.byType(AppButton), findsOneWidget);
      });
    });

    group('validation on empty submit', () {
      testWidgets('shows nickname validation error when nickname field is empty', (tester) async {
        await tester.pumpWidget(_buildRegisterScreen());
        await tester.pump();

        await tester.tap(find.byType(ElevatedButton));
        await tester.pump();

        expect(find.text('닉네임을 입력해주세요'), findsOneWidget);
      });

      testWidgets('shows email validation error when only nickname filled', (tester) async {
        await tester.pumpWidget(_buildRegisterScreen());
        await tester.pump();

        final fields = find.byType(TextFormField);
        await tester.enterText(fields.at(0), 'Alice');
        await tester.tap(find.byType(ElevatedButton));
        await tester.pump();

        expect(find.text('이메일을 입력해주세요'), findsOneWidget);
      });

      testWidgets('shows password validation error when only nickname and email filled', (tester) async {
        await tester.pumpWidget(_buildRegisterScreen());
        await tester.pump();

        final fields = find.byType(TextFormField);
        await tester.enterText(fields.at(0), 'Alice');
        await tester.enterText(fields.at(1), 'alice@example.com');
        await tester.tap(find.byType(ElevatedButton));
        await tester.pump();

        expect(find.text('비밀번호를 입력해주세요'), findsOneWidget);
      });

      testWidgets('shows password length error for password shorter than 8 chars', (tester) async {
        await tester.pumpWidget(_buildRegisterScreen());
        await tester.pump();

        final fields = find.byType(TextFormField);
        await tester.enterText(fields.at(0), 'Alice');
        await tester.enterText(fields.at(1), 'alice@example.com');
        await tester.enterText(fields.at(2), 'abc1');
        await tester.tap(find.byType(ElevatedButton));
        await tester.pump();

        expect(find.text('비밀번호는 8자 이상이어야 합니다'), findsOneWidget);
      });
    });

    group('loading state', () {
      testWidgets('button shows CircularProgressIndicator when status is loading', (tester) async {
        await tester.pumpWidget(
          _buildRegisterScreen(
            authState: const AuthState(status: AuthStatus.loading),
          ),
        );
        await tester.pump();

        expect(find.byType(CircularProgressIndicator), findsOneWidget);
      });
    });
  });
}
