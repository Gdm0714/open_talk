import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';

import 'package:open_talk/shared/widgets/avatar_widget.dart';

import '../../helpers/test_helpers.dart';

void main() {
  group('AvatarWidget', () {
    group('fallback initials', () {
      testWidgets('shows first character uppercased when no imageUrl is given', (tester) async {
        await tester.pumpWidget(
          buildTestableWidget(
            const AvatarWidget(name: 'alice'),
          ),
        );

        expect(find.text('A'), findsOneWidget);
      });

      testWidgets('shows two initials for a two-word name', (tester) async {
        await tester.pumpWidget(
          buildTestableWidget(
            const AvatarWidget(name: 'Hong Gildong'),
          ),
        );

        expect(find.text('HG'), findsOneWidget);
      });

      testWidgets('shows "?" when name is empty string', (tester) async {
        await tester.pumpWidget(
          buildTestableWidget(
            const AvatarWidget(name: ''),
          ),
        );

        expect(find.text('?'), findsOneWidget);
      });

      testWidgets('shows first character when name has only one word', (tester) async {
        await tester.pumpWidget(
          buildTestableWidget(
            const AvatarWidget(name: 'Bob'),
          ),
        );

        expect(find.text('B'), findsOneWidget);
      });
    });

    group('online indicator', () {
      testWidgets('does not render indicator container when showOnlineIndicator is false', (tester) async {
        await tester.pumpWidget(
          buildTestableWidget(
            const AvatarWidget(
              name: 'Alice',
              showOnlineIndicator: false,
            ),
          ),
        );

        // The indicator is a Positioned > Container — there should be no
        // Positioned widget when showOnlineIndicator is false.
        expect(find.byType(Positioned), findsNothing);
      });

      testWidgets('renders indicator when showOnlineIndicator is true', (tester) async {
        await tester.pumpWidget(
          buildTestableWidget(
            const AvatarWidget(
              name: 'Alice',
              showOnlineIndicator: true,
              isOnline: true,
            ),
          ),
        );

        expect(find.byType(Positioned), findsOneWidget);
      });

      testWidgets('renders indicator when isOnline is false but showOnlineIndicator is true', (tester) async {
        await tester.pumpWidget(
          buildTestableWidget(
            const AvatarWidget(
              name: 'Bob',
              showOnlineIndicator: true,
              isOnline: false,
            ),
          ),
        );

        expect(find.byType(Positioned), findsOneWidget);
      });
    });

    group('CircleAvatar', () {
      testWidgets('always renders a CircleAvatar', (tester) async {
        await tester.pumpWidget(
          buildTestableWidget(
            const AvatarWidget(name: 'Alice'),
          ),
        );

        expect(find.byType(CircleAvatar), findsOneWidget);
      });

      testWidgets('radius is half of the provided size', (tester) async {
        const size = 64.0;
        await tester.pumpWidget(
          buildTestableWidget(
            const AvatarWidget(name: 'Alice', size: size),
          ),
        );

        final circleAvatar = tester.widget<CircleAvatar>(find.byType(CircleAvatar));
        expect(circleAvatar.radius, size / 2);
      });
    });
  });
}
