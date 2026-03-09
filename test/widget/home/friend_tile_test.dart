import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';

import 'package:open_talk/features/home/data/models/friend_model.dart';
import 'package:open_talk/features/home/presentation/widgets/friend_tile.dart';
import 'package:open_talk/shared/widgets/avatar_widget.dart';

import '../../helpers/test_helpers.dart';

FriendModel _makeFriend({
  String nickname = 'Bob',
  String? statusMessage,
  String? avatarUrl,
}) {
  return FriendModel(
    id: 'friendship-1',
    friendId: 'user-bob',
    friendNickname: nickname,
    friendStatusMessage: statusMessage,
    friendAvatarUrl: avatarUrl,
  );
}

void main() {
  group('FriendTile', () {
    testWidgets('displays friend nickname', (tester) async {
      await tester.pumpWidget(
        buildTestableWidget(
          FriendTile(friend: _makeFriend(nickname: 'Charlie')),
        ),
      );

      expect(find.text('Charlie'), findsOneWidget);
    });

    testWidgets('displays status message when it is non-empty', (tester) async {
      await tester.pumpWidget(
        buildTestableWidget(
          FriendTile(
            friend: _makeFriend(statusMessage: '지금 바쁩니다'),
          ),
        ),
      );

      expect(find.text('지금 바쁩니다'), findsOneWidget);
    });

    testWidgets('does not display status message widget when statusMessage is null', (tester) async {
      await tester.pumpWidget(
        buildTestableWidget(
          FriendTile(friend: _makeFriend(statusMessage: null)),
        ),
      );

      // Nickname is visible
      expect(find.text('Bob'), findsOneWidget);
      // The FriendTile only renders a status-message Text when
      // friendStatusMessage is non-null and non-empty. With null status
      // there are exactly 2 Text descendants: the initials ('B') inside
      // AvatarWidget and the nickname ('Bob'). No third Text for status.
      final friendTile = find.byType(FriendTile);
      final textsInsideTile = tester
          .widgetList<Text>(
            find.descendant(of: friendTile, matching: find.byType(Text)),
          )
          .map((t) => t.data)
          .where((d) => d != null && d.isNotEmpty)
          .toList();
      // Should contain initials + nickname but NO status message string
      expect(textsInsideTile.length, 2);
      expect(textsInsideTile, contains('Bob'));
    });

    testWidgets('does not display status message when it is empty string', (tester) async {
      await tester.pumpWidget(
        buildTestableWidget(
          FriendTile(friend: _makeFriend(statusMessage: '')),
        ),
      );

      // Empty status message should not render a Text widget for it
      expect(find.text(''), findsNothing);
    });

    testWidgets('renders AvatarWidget with friend nickname as initials source', (tester) async {
      await tester.pumpWidget(
        buildTestableWidget(
          FriendTile(friend: _makeFriend(nickname: 'Dave')),
        ),
      );

      final avatar = tester.widget<AvatarWidget>(find.byType(AvatarWidget));
      expect(avatar.name, 'Dave');
    });

    testWidgets('avatar shows initials when no avatarUrl is provided', (tester) async {
      await tester.pumpWidget(
        buildTestableWidget(
          FriendTile(friend: _makeFriend(nickname: 'Eve', avatarUrl: null)),
        ),
      );

      // First character of 'Eve' uppercased
      expect(find.text('E'), findsOneWidget);
    });

    testWidgets('calls onTap callback when tile is tapped', (tester) async {
      var tapped = false;
      await tester.pumpWidget(
        buildTestableWidget(
          FriendTile(
            friend: _makeFriend(),
            onTap: () => tapped = true,
          ),
        ),
      );

      await tester.tap(find.byType(InkWell));
      await tester.pump();

      expect(tapped, isTrue);
    });
  });
}
