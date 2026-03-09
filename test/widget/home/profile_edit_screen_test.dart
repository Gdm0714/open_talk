import 'package:dio/dio.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:mocktail/mocktail.dart';

import 'package:open_talk/core/network/api_client.dart';
import 'package:open_talk/features/auth/data/models/user_model.dart';
import 'package:open_talk/features/auth/domain/providers/auth_provider.dart';
import 'package:open_talk/features/home/presentation/screens/profile_edit_screen.dart';

class MockDio extends Mock implements Dio {}

class MockAuthNotifier extends StateNotifier<AuthState>
    implements AuthNotifier {
  MockAuthNotifier(AuthState state) : super(state);

  @override
  Future<void> checkAuthStatus() async {}

  @override
  Future<void> login({required String email, required String password}) async {}

  @override
  Future<void> register({
    required String email,
    required String password,
    required String nickname,
  }) async {}

  @override
  Future<void> logout() async {}

  @override
  void clearError() {}
}

void main() {
  late MockDio mockDio;

  setUp(() {
    mockDio = MockDio();
    registerFallbackValue(
      Options(),
    );
  });

  const testUser = UserModel(
    id: 'user-1',
    email: 'test@example.com',
    nickname: 'TestUser',
    statusMessage: 'Hello world',
  );

  Widget buildSubject({UserModel? user}) {
    final authState = AuthState(
      status: AuthStatus.authenticated,
      user: user ?? testUser,
    );

    return ProviderScope(
      overrides: [
        apiClientProvider.overrideWithValue(mockDio),
        authStateProvider.overrideWith(
          (ref) => MockAuthNotifier(authState),
        ),
      ],
      child: const MaterialApp(
        home: ProfileEditScreen(),
      ),
    );
  }

  testWidgets('renders nickname and status message fields', (tester) async {
    await tester.pumpWidget(buildSubject());
    await tester.pumpAndSettle();

    expect(find.text('닉네임'), findsOneWidget);
    expect(find.text('상태 메시지'), findsOneWidget);
  });

  testWidgets('populates initial values from provider', (tester) async {
    await tester.pumpWidget(buildSubject());
    await tester.pumpAndSettle();

    final nicknameField = find.widgetWithText(TextFormField, 'TestUser');
    expect(nicknameField, findsOneWidget);

    final statusField = find.widgetWithText(TextFormField, 'Hello world');
    expect(statusField, findsOneWidget);
  });

  testWidgets('save button is rendered', (tester) async {
    await tester.pumpWidget(buildSubject());
    await tester.pumpAndSettle();

    expect(find.text('저장'), findsOneWidget);
  });

  testWidgets('empty nickname shows validation error', (tester) async {
    await tester.pumpWidget(buildSubject());
    await tester.pumpAndSettle();

    // Clear the nickname field
    final nicknameField = find.widgetWithText(TextFormField, 'TestUser');
    await tester.tap(nicknameField);
    await tester.pump();
    await tester.enterText(nicknameField, '');

    // Tap save to trigger validation
    await tester.tap(find.text('저장'));
    await tester.pump();

    expect(find.text('닉네임을 입력하세요'), findsOneWidget);
  });

  testWidgets('short nickname shows length validation error', (tester) async {
    await tester.pumpWidget(buildSubject());
    await tester.pumpAndSettle();

    final nicknameField = find.widgetWithText(TextFormField, 'TestUser');
    await tester.tap(nicknameField);
    await tester.pump();
    await tester.enterText(nicknameField, 'A');

    await tester.tap(find.text('저장'));
    await tester.pump();

    expect(find.text('닉네임은 2자 이상이어야 합니다'), findsOneWidget);
  });

  testWidgets('save button is tappable when form is valid', (tester) async {
    when(() => mockDio.put(any(), data: any(named: 'data')))
        .thenAnswer((_) async => Response(
              requestOptions: RequestOptions(path: '/users/me'),
              statusCode: 200,
            ));

    await tester.pumpWidget(buildSubject());
    await tester.pumpAndSettle();

    final saveButton = find.text('저장');
    expect(saveButton, findsOneWidget);
    await tester.tap(saveButton);
    await tester.pump();
    // No crash — button was tappable
  });
}
