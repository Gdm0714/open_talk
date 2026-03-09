import 'package:flutter_test/flutter_test.dart';

import 'package:open_talk/features/chat/data/models/chat_room_model.dart';
import 'package:open_talk/features/auth/data/models/user_model.dart';

void main() {
  group('ChatRoomModel', () {
    final memberAlice = {
      'id': 'user-alice',
      'email': 'alice@example.com',
      'nickname': 'Alice',
      'avatar_url': null,
      'status_message': null,
    };
    final memberBob = {
      'id': 'user-bob',
      'email': 'bob@example.com',
      'nickname': 'Bob',
      'avatar_url': null,
      'status_message': null,
    };

    group('fromJson', () {
      test('parses direct chat room with all fields', () {
        final json = {
          'id': 'room-1',
          'name': null,
          'type': 'direct',
          'last_message': 'Hey!',
          'last_message_at': '2024-01-15T10:30:00.000Z',
          'members': [memberAlice, memberBob],
          'unread_count': 3,
        };

        final model = ChatRoomModel.fromJson(json);

        expect(model.id, 'room-1');
        expect(model.name, isNull);
        expect(model.type, ChatRoomType.direct);
        expect(model.lastMessage, 'Hey!');
        expect(model.lastMessageAt, isNotNull);
        expect(model.members.length, 2);
        expect(model.unreadCount, 3);
      });

      test('parses group chat room type correctly', () {
        final json = {
          'id': 'room-2',
          'name': 'Team Chat',
          'type': 'group',
          'last_message': null,
          'last_message_at': null,
          'members': [],
          'unread_count': 0,
        };

        final model = ChatRoomModel.fromJson(json);

        expect(model.type, ChatRoomType.group);
      });

      test('defaults type to direct when value is not group', () {
        final json = {
          'id': 'room-3',
          'name': null,
          'type': 'unknown_value',
          'last_message': null,
          'last_message_at': null,
          'members': [],
          'unread_count': 0,
        };

        final model = ChatRoomModel.fromJson(json);

        expect(model.type, ChatRoomType.direct);
      });

      test('parses lastMessageAt as DateTime from ISO 8601 string', () {
        final json = {
          'id': 'room-4',
          'name': null,
          'type': 'direct',
          'last_message': 'hi',
          'last_message_at': '2024-06-01T08:00:00.000Z',
          'members': [],
          'unread_count': 0,
        };

        final model = ChatRoomModel.fromJson(json);

        expect(model.lastMessageAt, DateTime.parse('2024-06-01T08:00:00.000Z'));
      });

      test('sets lastMessageAt to null when json value is null', () {
        final json = {
          'id': 'room-5',
          'name': null,
          'type': 'direct',
          'last_message': null,
          'last_message_at': null,
          'members': [],
          'unread_count': 0,
        };

        final model = ChatRoomModel.fromJson(json);

        expect(model.lastMessageAt, isNull);
      });

      test('defaults unreadCount to 0 when key is absent', () {
        final json = {
          'id': 'room-6',
          'name': null,
          'type': 'direct',
          'last_message': null,
          'last_message_at': null,
          'members': [],
        };

        final model = ChatRoomModel.fromJson(json);

        expect(model.unreadCount, 0);
      });

      test('parses members list into UserModel instances', () {
        final json = {
          'id': 'room-7',
          'name': null,
          'type': 'direct',
          'last_message': null,
          'last_message_at': null,
          'members': [memberAlice, memberBob],
          'unread_count': 0,
        };

        final model = ChatRoomModel.fromJson(json);

        expect(model.members[0], isA<UserModel>());
        expect(model.members[0].nickname, 'Alice');
        expect(model.members[1].nickname, 'Bob');
      });
    });

    group('toJson', () {
      test('serializes type direct as string "direct"', () {
        const model = ChatRoomModel(
          id: 'room-1',
          type: ChatRoomType.direct,
        );

        final json = model.toJson();

        expect(json['type'], 'direct');
      });

      test('serializes type group as string "group"', () {
        const model = ChatRoomModel(
          id: 'room-2',
          type: ChatRoomType.group,
          name: 'Friends',
        );

        final json = model.toJson();

        expect(json['type'], 'group');
      });

      test('serializes lastMessageAt as ISO 8601 string', () {
        final dt = DateTime.utc(2024, 6, 1, 8, 0, 0);
        final model = ChatRoomModel(
          id: 'room-3',
          type: ChatRoomType.direct,
          lastMessageAt: dt,
        );

        final json = model.toJson();

        expect(json['last_message_at'], dt.toIso8601String());
      });

      test('serializes members as list of maps', () {
        const alice = UserModel(
          id: 'user-alice',
          email: 'alice@example.com',
          nickname: 'Alice',
        );
        const model = ChatRoomModel(
          id: 'room-4',
          type: ChatRoomType.direct,
          members: [alice],
        );

        final json = model.toJson();

        expect(json['members'], isA<List>());
        expect((json['members'] as List).first, isA<Map<String, dynamic>>());
      });
    });

    group('displayName', () {
      const currentUserId = 'user-alice';

      test('returns name field when it is non-empty', () {
        const model = ChatRoomModel(
          id: 'room-1',
          name: 'Project Alpha',
          type: ChatRoomType.group,
        );

        expect(model.displayName(currentUserId), 'Project Alpha');
      });

      test('returns other member nickname for direct chat without name', () {
        const alice = UserModel(
          id: 'user-alice',
          email: 'alice@example.com',
          nickname: 'Alice',
        );
        const bob = UserModel(
          id: 'user-bob',
          email: 'bob@example.com',
          nickname: 'Bob',
        );
        const model = ChatRoomModel(
          id: 'room-2',
          type: ChatRoomType.direct,
          members: [alice, bob],
        );

        expect(model.displayName(currentUserId), 'Bob');
      });

      test('returns fallback text when direct chat has no other member', () {
        const model = ChatRoomModel(
          id: 'room-3',
          type: ChatRoomType.direct,
          members: [],
        );

        expect(model.displayName(currentUserId), '알 수 없는 사용자');
      });

      test('returns joined member nicknames for group chat without name', () {
        const alice = UserModel(
          id: 'user-alice',
          email: 'alice@example.com',
          nickname: 'Alice',
        );
        const bob = UserModel(
          id: 'user-bob',
          email: 'bob@example.com',
          nickname: 'Bob',
        );
        const model = ChatRoomModel(
          id: 'room-4',
          type: ChatRoomType.group,
          members: [alice, bob],
        );

        expect(model.displayName(currentUserId), 'Alice, Bob');
      });
    });
  });
}
