import 'package:flutter_test/flutter_test.dart';

import 'package:open_talk/features/chat/data/models/message_model.dart';

void main() {
  group('MessageModel', () {
    final validJson = {
      'id': 'msg-1',
      'chat_room_id': 'room-1',
      'sender_id': 'user-alice',
      'sender_nickname': 'Alice',
      'content': 'Hello world',
      'message_type': 'text',
      'created_at': '2024-01-15T10:30:00.000Z',
    };

    group('fromJson', () {
      test('parses all fields from valid json', () {
        final model = MessageModel.fromJson(validJson);

        expect(model.id, 'msg-1');
        expect(model.chatRoomId, 'room-1');
        expect(model.senderId, 'user-alice');
        expect(model.senderNickname, 'Alice');
        expect(model.content, 'Hello world');
        expect(model.messageType, MessageType.text);
      });

      test('parses createdAt as DateTime from ISO 8601 string', () {
        final model = MessageModel.fromJson(validJson);

        expect(
          model.createdAt,
          DateTime.parse('2024-01-15T10:30:00.000Z'),
        );
      });

      test('parses message_type "image" as MessageType.image', () {
        final json = {...validJson, 'message_type': 'image'};

        final model = MessageModel.fromJson(json);

        expect(model.messageType, MessageType.image);
      });

      test('parses message_type "system" as MessageType.system', () {
        final json = {...validJson, 'message_type': 'system'};

        final model = MessageModel.fromJson(json);

        expect(model.messageType, MessageType.system);
      });

      test('defaults to MessageType.text when message_type is null', () {
        final json = {...validJson, 'message_type': null};

        final model = MessageModel.fromJson(json);

        expect(model.messageType, MessageType.text);
      });

      test('defaults to MessageType.text when message_type is unknown string', () {
        final json = {...validJson, 'message_type': 'video'};

        final model = MessageModel.fromJson(json);

        expect(model.messageType, MessageType.text);
      });

      test('defaults senderNickname to empty string when key is null', () {
        final json = {...validJson, 'sender_nickname': null};

        final model = MessageModel.fromJson(json);

        expect(model.senderNickname, '');
      });
    });

    group('toJson', () {
      test('serializes all fields correctly', () {
        final createdAt = DateTime.utc(2024, 1, 15, 10, 30, 0);
        final model = MessageModel(
          id: 'msg-1',
          chatRoomId: 'room-1',
          senderId: 'user-alice',
          senderNickname: 'Alice',
          content: 'Hello world',
          messageType: MessageType.text,
          createdAt: createdAt,
        );

        final json = model.toJson();

        expect(json['id'], 'msg-1');
        expect(json['chat_room_id'], 'room-1');
        expect(json['sender_id'], 'user-alice');
        expect(json['sender_nickname'], 'Alice');
        expect(json['content'], 'Hello world');
        expect(json['message_type'], 'text');
        expect(json['created_at'], createdAt.toIso8601String());
      });

      test('serializes MessageType.image as string "image"', () {
        final model = MessageModel(
          id: 'msg-2',
          chatRoomId: 'room-1',
          senderId: 'user-alice',
          senderNickname: 'Alice',
          content: 'image data',
          messageType: MessageType.image,
          createdAt: DateTime.now(),
        );

        expect(model.toJson()['message_type'], 'image');
      });

      test('serializes MessageType.system as string "system"', () {
        final model = MessageModel(
          id: 'msg-3',
          chatRoomId: 'room-1',
          senderId: 'system',
          senderNickname: '',
          content: 'Alice joined the chat',
          messageType: MessageType.system,
          createdAt: DateTime.now(),
        );

        expect(model.toJson()['message_type'], 'system');
      });
    });

    group('isMine', () {
      test('returns true when senderId matches currentUserId', () {
        final model = MessageModel(
          id: 'msg-1',
          chatRoomId: 'room-1',
          senderId: 'user-alice',
          senderNickname: 'Alice',
          content: 'Hi',
          createdAt: DateTime.now(),
        );

        expect(model.isMine('user-alice'), isTrue);
      });

      test('returns false when senderId does not match currentUserId', () {
        final model = MessageModel(
          id: 'msg-2',
          chatRoomId: 'room-1',
          senderId: 'user-bob',
          senderNickname: 'Bob',
          content: 'Hey',
          createdAt: DateTime.now(),
        );

        expect(model.isMine('user-alice'), isFalse);
      });
    });
  });
}
