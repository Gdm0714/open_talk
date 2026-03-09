import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';

import 'package:open_talk/shared/widgets/app_button.dart';

import '../../helpers/test_helpers.dart';

void main() {
  group('AppButton', () {
    group('primary variant', () {
      testWidgets('renders as ElevatedButton with correct label', (tester) async {
        await tester.pumpWidget(
          buildTestableWidget(
            const AppButton(text: '확인'),
          ),
        );

        expect(find.byType(ElevatedButton), findsOneWidget);
        expect(find.text('확인'), findsOneWidget);
      });

      testWidgets('calls onPressed callback when tapped', (tester) async {
        var tapped = false;
        await tester.pumpWidget(
          buildTestableWidget(
            AppButton(
              text: '클릭',
              onPressed: () => tapped = true,
            ),
          ),
        );

        await tester.tap(find.byType(ElevatedButton));
        await tester.pump();

        expect(tapped, isTrue);
      });

      testWidgets('does not fire onPressed when enabled is false', (tester) async {
        var tapped = false;
        await tester.pumpWidget(
          buildTestableWidget(
            AppButton(
              text: '비활성',
              onPressed: () => tapped = true,
              enabled: false,
            ),
          ),
        );

        await tester.tap(find.byType(ElevatedButton));
        await tester.pump();

        expect(tapped, isFalse);
      });
    });

    group('outline variant', () {
      testWidgets('renders as OutlinedButton with correct label', (tester) async {
        await tester.pumpWidget(
          buildTestableWidget(
            const AppButton(
              text: '취소',
              variant: AppButtonVariant.outline,
            ),
          ),
        );

        expect(find.byType(OutlinedButton), findsOneWidget);
        expect(find.text('취소'), findsOneWidget);
      });

      testWidgets('calls onPressed when outline button is tapped', (tester) async {
        var tapped = false;
        await tester.pumpWidget(
          buildTestableWidget(
            AppButton(
              text: '취소',
              variant: AppButtonVariant.outline,
              onPressed: () => tapped = true,
            ),
          ),
        );

        await tester.tap(find.byType(OutlinedButton));
        await tester.pump();

        expect(tapped, isTrue);
      });
    });

    group('loading state', () {
      testWidgets('shows CircularProgressIndicator when isLoading is true', (tester) async {
        await tester.pumpWidget(
          buildTestableWidget(
            const AppButton(
              text: '로딩중',
              isLoading: true,
            ),
          ),
        );

        expect(find.byType(CircularProgressIndicator), findsOneWidget);
        // Label text should not be visible while loading
        expect(find.text('로딩중'), findsNothing);
      });

      testWidgets('does not fire onPressed while isLoading is true', (tester) async {
        var tapped = false;
        await tester.pumpWidget(
          buildTestableWidget(
            AppButton(
              text: '로딩중',
              isLoading: true,
              onPressed: () => tapped = true,
            ),
          ),
        );

        await tester.tap(find.byType(ElevatedButton));
        await tester.pump();

        expect(tapped, isFalse);
      });
    });

    group('with icon', () {
      testWidgets('renders icon alongside text when icon is provided', (tester) async {
        await tester.pumpWidget(
          buildTestableWidget(
            const AppButton(
              text: '아이콘 버튼',
              icon: Icons.add,
            ),
          ),
        );

        expect(find.byIcon(Icons.add), findsOneWidget);
        expect(find.text('아이콘 버튼'), findsOneWidget);
      });
    });
  });
}
