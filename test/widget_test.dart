import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import 'package:open_talk/features/auth/presentation/screens/login_screen.dart';
import 'package:open_talk/core/theme/app_theme.dart';

void main() {
  testWidgets('LoginScreen renders correctly', (WidgetTester tester) async {
    await tester.pumpWidget(
      ProviderScope(
        child: MaterialApp(
          theme: AppTheme.light,
          home: const LoginScreen(),
        ),
      ),
    );

    expect(find.text('OpenTalk'), findsOneWidget);
    expect(find.text('로그인'), findsWidgets);
    expect(find.text('회원가입'), findsOneWidget);
  });
}
