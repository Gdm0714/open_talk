import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:mocktail/mocktail.dart';
import 'package:dio/dio.dart';

import 'package:open_talk/core/network/api_client.dart';
import 'package:open_talk/features/home/presentation/screens/password_change_screen.dart';

class MockDio extends Mock implements Dio {}

void main() {
  late MockDio mockDio;

  setUp(() {
    mockDio = MockDio();
    registerFallbackValue(Options());
  });

  Widget buildSubject() {
    return ProviderScope(
      overrides: [
        apiClientProvider.overrideWithValue(mockDio),
      ],
      child: const MaterialApp(
        home: PasswordChangeScreen(),
      ),
    );
  }

  testWidgets('renders all three password fields', (tester) async {
    await tester.pumpWidget(buildSubject());
    await tester.pumpAndSettle();

    expect(find.text('현재 비밀번호'), findsOneWidget);
    expect(find.text('새 비밀번호'), findsOneWidget);
    expect(find.text('새 비밀번호 확인'), findsOneWidget);
  });

  testWidgets('change button is rendered', (tester) async {
    await tester.pumpWidget(buildSubject());
    await tester.pumpAndSettle();

    expect(find.text('변경'), findsOneWidget);
  });

  testWidgets('empty current password shows validation error', (tester) async {
    await tester.pumpWidget(buildSubject());
    await tester.pumpAndSettle();

    await tester.tap(find.text('변경'));
    await tester.pump();

    expect(find.text('현재 비밀번호를 입력하세요'), findsOneWidget);
  });

  testWidgets('empty new password shows validation error', (tester) async {
    await tester.pumpWidget(buildSubject());
    await tester.pumpAndSettle();

    // Fill current password but leave new password empty
    final currentPwField = find.widgetWithText(TextFormField, '');
    await tester.enterText(currentPwField.first, 'oldpassword');

    await tester.tap(find.text('변경'));
    await tester.pump();

    expect(find.text('새 비밀번호를 입력하세요'), findsOneWidget);
  });

  testWidgets('short new password shows length validation error',
      (tester) async {
    await tester.pumpWidget(buildSubject());
    await tester.pumpAndSettle();

    final fields = find.byType(TextFormField);
    await tester.enterText(fields.at(0), 'oldpassword');
    await tester.enterText(fields.at(1), 'short');

    await tester.tap(find.text('변경'));
    await tester.pump();

    expect(find.text('비밀번호는 8자 이상이어야 합니다'), findsOneWidget);
  });

  testWidgets('mismatched confirm password shows validation error',
      (tester) async {
    await tester.pumpWidget(buildSubject());
    await tester.pumpAndSettle();

    final fields = find.byType(TextFormField);
    await tester.enterText(fields.at(0), 'oldpassword');
    await tester.enterText(fields.at(1), 'newpassword1');
    await tester.enterText(fields.at(2), 'newpassword2');

    await tester.tap(find.text('변경'));
    await tester.pump();

    expect(find.text('비밀번호가 일치하지 않습니다'), findsOneWidget);
  });

  testWidgets('empty confirm password shows validation error', (tester) async {
    await tester.pumpWidget(buildSubject());
    await tester.pumpAndSettle();

    final fields = find.byType(TextFormField);
    await tester.enterText(fields.at(0), 'oldpassword');
    await tester.enterText(fields.at(1), 'newpassword1');
    // leave confirm empty

    await tester.tap(find.text('변경'));
    await tester.pump();

    expect(find.text('비밀번호를 다시 입력하세요'), findsOneWidget);
  });
}
