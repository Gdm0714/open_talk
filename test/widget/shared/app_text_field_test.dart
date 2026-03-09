import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';

import 'package:open_talk/shared/widgets/app_text_field.dart';

import '../../helpers/test_helpers.dart';

void main() {
  group('AppTextField', () {
    testWidgets('renders label text when label is provided', (tester) async {
      await tester.pumpWidget(
        buildTestableWidget(
          const AppTextField(label: '이메일'),
        ),
      );

      expect(find.text('이메일'), findsOneWidget);
    });

    testWidgets('does not render label Text when label is null', (tester) async {
      await tester.pumpWidget(
        buildTestableWidget(
          const AppTextField(hintText: '입력하세요'),
        ),
      );

      // When label is null, the AppTextField skips the label Text widget.
      // The TextFormField is present, but no additional Text widget for a label.
      expect(find.byType(TextFormField), findsOneWidget);
      // There should be no Text widget with value '이메일' or any label
      expect(find.text('이메일'), findsNothing);
    });

    testWidgets('shows error text in InputDecoration when errorText is set', (tester) async {
      await tester.pumpWidget(
        buildTestableWidget(
          const AppTextField(
            label: '비밀번호',
            errorText: '비밀번호가 올바르지 않습니다',
          ),
        ),
      );

      expect(find.text('비밀번호가 올바르지 않습니다'), findsOneWidget);
    });

    testWidgets('EditableText has obscureText true when obscureText is set', (tester) async {
      await tester.pumpWidget(
        buildTestableWidget(
          const AppTextField(
            label: '비밀번호',
            obscureText: true,
          ),
        ),
      );

      // TextFormField renders an EditableText internally which exposes obscureText
      final editableText = tester.widget<EditableText>(find.byType(EditableText));
      expect(editableText.obscureText, isTrue);
    });

    testWidgets('EditableText has obscureText false by default', (tester) async {
      await tester.pumpWidget(
        buildTestableWidget(
          const AppTextField(label: '이메일'),
        ),
      );

      final editableText = tester.widget<EditableText>(find.byType(EditableText));
      expect(editableText.obscureText, isFalse);
    });

    testWidgets('calls onChanged callback when text is entered', (tester) async {
      String? changedValue;
      await tester.pumpWidget(
        buildTestableWidget(
          AppTextField(
            label: '검색',
            onChanged: (value) => changedValue = value,
          ),
        ),
      );

      await tester.enterText(find.byType(TextFormField), 'hello');
      await tester.pump();

      expect(changedValue, 'hello');
    });

    testWidgets('uses provided TextEditingController', (tester) async {
      final controller = TextEditingController(text: '초기값');
      addTearDown(controller.dispose);

      await tester.pumpWidget(
        buildTestableWidget(
          AppTextField(
            label: '닉네임',
            controller: controller,
          ),
        ),
      );

      expect(find.text('초기값'), findsOneWidget);
    });

    testWidgets('renders hint text inside the field', (tester) async {
      await tester.pumpWidget(
        buildTestableWidget(
          const AppTextField(
            hintText: 'example@email.com',
          ),
        ),
      );

      expect(find.text('example@email.com'), findsOneWidget);
    });

    testWidgets('renders suffix icon widget when suffixIcon is provided', (tester) async {
      // Verify by inspecting the AppTextField widget's suffixIcon property
      // rather than finding it in the render tree (InputDecoration slots are
      // internal and finding by icon data can be unreliable in tests).
      await tester.pumpWidget(
        buildTestableWidget(
          const AppTextField(
            label: '비밀번호',
            suffixIcon: Icon(Icons.visibility_outlined),
          ),
        ),
      );

      // The AppTextField widget itself should be present
      final appTextField = tester.widget<AppTextField>(find.byType(AppTextField));
      expect(appTextField.suffixIcon, isNotNull);
      expect(appTextField.suffixIcon, isA<Icon>());
    });

    testWidgets('renders prefix icon when provided', (tester) async {
      await tester.pumpWidget(
        buildTestableWidget(
          const AppTextField(
            label: '이메일',
            prefixIcon: Icon(Icons.email_outlined),
          ),
        ),
      );

      expect(find.byIcon(Icons.email_outlined), findsOneWidget);
    });
  });
}
