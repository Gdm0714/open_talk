import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:intl/date_symbol_data_local.dart';

import 'package:open_talk/features/auth/data/models/user_model.dart';
import 'package:open_talk/features/chat/data/models/chat_room_model.dart';
import 'package:open_talk/features/chat/presentation/widgets/chat_room_tile.dart';

import '../../helpers/test_helpers.dart';

const _currentUserId = 'user-alice';

const _alice = UserModel(
  id: _currentUserId,
  email: 'alice@example.com',
  nickname: 'Alice',
);

const _bob = UserModel(
  id: 'user-bob',
  email: 'bob@example.com',
  nickname: 'Bob',
);

ChatRoomModel _directRoom({
  String? lastMessage,
  DateTime? lastMessageAt,
  int unreadCount = 0,
}) {
  return ChatRoomModel(
    id: 'room-1',
    type: ChatRoomType.direct,
    members: const [_alice, _bob],
    lastMessage: lastMessage,
    lastMessageAt: lastMessageAt,
    unreadCount: unreadCount,
  );
}

void main() {
  setUpAll(() async {
    await initializeDateFormatting('ko_KR', null);
  });

  group('ChatRoomTile', () {
    testWidgets('displays the room display name', (tester) async {
      await tester.pumpWidget(
        buildTestableWidget(
          ChatRoomTile(
            chatRoom: _directRoom(),
            currentUserId: _currentUserId,
          ),
        ),
      );

      // For a direct room without a name, displayName returns the other
      // member's nickname, which is 'Bob'.
      expect(find.text('Bob'), findsOneWidget);
    });

    testWidgets('displays named group room name', (tester) async {
      const room = ChatRoomModel(
        id: 'room-2',
        name: 'Team Alpha',
        type: ChatRoomType.group,
        members: [_alice, _bob],
      );

      await tester.pumpWidget(
        buildTestableWidget(
          ChatRoomTile(
            chatRoom: room,
            currentUserId: _currentUserId,
          ),
        ),
      );

      expect(find.text('Team Alpha'), findsOneWidget);
    });

    testWidgets('displays last message preview text', (tester) async {
      await tester.pumpWidget(
        buildTestableWidget(
          ChatRoomTile(
            chatRoom: _directRoom(lastMessage: '안녕하세요!'),
            currentUserId: _currentUserId,
          ),
        ),
      );

      expect(find.text('안녕하세요!'), findsOneWidget);
    });

    testWidgets('shows empty string when lastMessage is null', (tester) async {
      await tester.pumpWidget(
        buildTestableWidget(
          ChatRoomTile(
            chatRoom: _directRoom(lastMessage: null),
            currentUserId: _currentUserId,
          ),
        ),
      );

      // Should not throw; the tile still renders
      expect(find.byType(ChatRoomTile), findsOneWidget);
    });

    testWidgets('displays unread badge when unreadCount is greater than 0', (tester) async {
      await tester.pumpWidget(
        buildTestableWidget(
          ChatRoomTile(
            chatRoom: _directRoom(unreadCount: 5),
            currentUserId: _currentUserId,
          ),
        ),
      );

      expect(find.text('5'), findsOneWidget);
    });

    testWidgets('does not display unread badge when unreadCount is 0', (tester) async {
      await tester.pumpWidget(
        buildTestableWidget(
          ChatRoomTile(
            chatRoom: _directRoom(unreadCount: 0),
            currentUserId: _currentUserId,
          ),
        ),
      );

      // The badge container only exists when unreadCount > 0
      // '0' should not appear as badge text
      expect(find.text('0'), findsNothing);
    });

    testWidgets('shows "99+" badge text when unreadCount exceeds 99', (tester) async {
      await tester.pumpWidget(
        buildTestableWidget(
          ChatRoomTile(
            chatRoom: _directRoom(unreadCount: 150),
            currentUserId: _currentUserId,
          ),
        ),
      );

      expect(find.text('99+'), findsOneWidget);
    });

    testWidgets('calls onTap callback when tile is tapped', (tester) async {
      var tapped = false;
      await tester.pumpWidget(
        buildTestableWidget(
          ChatRoomTile(
            chatRoom: _directRoom(),
            currentUserId: _currentUserId,
            onTap: () => tapped = true,
          ),
        ),
      );

      await tester.tap(find.byType(InkWell));
      await tester.pump();

      expect(tapped, isTrue);
    });

    testWidgets('displays formatted last message time when lastMessageAt is set', (tester) async {
      // Use a fixed date 1 day before today at midnight to guarantee '어제'
      final now = DateTime.now();
      final yesterday = DateTime(now.year, now.month, now.day - 1, 10, 0);
      await tester.pumpWidget(
        buildTestableWidget(
          ChatRoomTile(
            chatRoom: _directRoom(
              lastMessage: 'Hi',
              lastMessageAt: yesterday,
            ),
            currentUserId: _currentUserId,
          ),
        ),
      );
      await tester.pump();

      expect(find.text('어제'), findsOneWidget);
    });
  });
}
