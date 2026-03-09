import 'package:flutter_test/flutter_test.dart';

import 'package:open_talk/shared/services/websocket_service.dart';

void main() {
  group('WebSocketMessage', () {
    group('fromJson', () {
      test('parses required type field', () {
        final msg = WebSocketMessage.fromJson({'type': 'message'});
        expect(msg.type, 'message');
      });

      test('parses room_id field', () {
        final msg = WebSocketMessage.fromJson({
          'type': 'message',
          'room_id': 'room-123',
        });
        expect(msg.roomId, 'room-123');
      });

      test('parses content field', () {
        final msg = WebSocketMessage.fromJson({
          'type': 'message',
          'content': 'Hello!',
        });
        expect(msg.content, 'Hello!');
      });

      test('parses sender_id field', () {
        final msg = WebSocketMessage.fromJson({
          'type': 'message',
          'sender_id': 'user-42',
        });
        expect(msg.senderId, 'user-42');
      });

      test('parses sender_name field', () {
        final msg = WebSocketMessage.fromJson({
          'type': 'message',
          'sender_name': 'Alice',
        });
        expect(msg.senderName, 'Alice');
      });

      test('stores entire json as data', () {
        final json = {
          'type': 'message',
          'room_id': 'room-1',
          'content': 'Hi',
        };
        final msg = WebSocketMessage.fromJson(json);
        expect(msg.data, json);
      });

      test('optional fields are null when absent', () {
        final msg = WebSocketMessage.fromJson({'type': 'ping'});
        expect(msg.roomId, isNull);
        expect(msg.content, isNull);
        expect(msg.senderId, isNull);
        expect(msg.senderName, isNull);
      });

      test('parses all fields together', () {
        final msg = WebSocketMessage.fromJson({
          'type': 'message',
          'room_id': 'room-99',
          'content': 'Hello world',
          'sender_id': 'user-1',
          'sender_name': 'Bob',
        });
        expect(msg.type, 'message');
        expect(msg.roomId, 'room-99');
        expect(msg.content, 'Hello world');
        expect(msg.senderId, 'user-1');
        expect(msg.senderName, 'Bob');
      });
    });

    group('toJson', () {
      test('serializes type', () {
        final msg = WebSocketMessage(type: 'ping');
        expect(msg.toJson()['type'], 'ping');
      });

      test('omits null optional fields', () {
        final msg = WebSocketMessage(type: 'ping');
        final json = msg.toJson();
        expect(json.containsKey('room_id'), isFalse);
        expect(json.containsKey('content'), isFalse);
        expect(json.containsKey('sender_id'), isFalse);
        expect(json.containsKey('sender_name'), isFalse);
      });

      test('includes room_id when set', () {
        final msg = WebSocketMessage(type: 'join', roomId: 'room-5');
        expect(msg.toJson()['room_id'], 'room-5');
      });

      test('includes content when set', () {
        final msg = WebSocketMessage(type: 'message', content: 'Hi there');
        expect(msg.toJson()['content'], 'Hi there');
      });

      test('includes sender_id when set', () {
        final msg = WebSocketMessage(type: 'message', senderId: 'user-7');
        expect(msg.toJson()['sender_id'], 'user-7');
      });

      test('includes sender_name when set', () {
        final msg = WebSocketMessage(type: 'message', senderName: 'Charlie');
        expect(msg.toJson()['sender_name'], 'Charlie');
      });

      test('round-trips through toJson and fromJson', () {
        final original = WebSocketMessage(
          type: 'message',
          roomId: 'room-1',
          content: 'Test',
          senderId: 'user-1',
          senderName: 'Alice',
        );
        final json = original.toJson();
        final restored = WebSocketMessage.fromJson(json);
        expect(restored.type, original.type);
        expect(restored.roomId, original.roomId);
        expect(restored.content, original.content);
        expect(restored.senderId, original.senderId);
        expect(restored.senderName, original.senderName);
      });
    });
  });

  group('WebSocketService', () {
    late WebSocketService service;

    setUp(() {
      service = WebSocketService();
    });

    tearDown(() {
      service.dispose();
    });

    test('initial state is not connected', () {
      expect(service.isConnected, isFalse);
    });

    test('exposes messageStream', () {
      expect(service.messageStream, isNotNull);
    });

    test('exposes connectionStream', () {
      expect(service.connectionStream, isNotNull);
    });

    test('disconnect sets isConnected to false', () {
      service.disconnect();
      expect(service.isConnected, isFalse);
    });

    test('sendRaw does nothing when not connected', () {
      // Should not throw when disconnected
      expect(
        () => service.sendRaw({'type': 'ping'}),
        returnsNormally,
      );
    });

    test('joinRoom does not throw when not connected', () {
      expect(() => service.joinRoom('room-1'), returnsNormally);
    });

    test('leaveRoom does not throw when not connected', () {
      expect(() => service.leaveRoom('room-1'), returnsNormally);
    });

    test('sendMessage does not throw when not connected', () {
      expect(
        () => service.sendMessage('room-1', 'hello'),
        returnsNormally,
      );
    });

    test('sendTyping does not throw when not connected', () {
      expect(() => service.sendTyping('room-1'), returnsNormally);
    });
  });
}
