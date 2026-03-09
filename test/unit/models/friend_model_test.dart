import 'package:flutter_test/flutter_test.dart';

import 'package:open_talk/features/home/data/models/friend_model.dart';

void main() {
  group('FriendModel', () {
    final validJson = {
      'id': 'friendship-1',
      'friend_id': 'user-bob',
      'friend_nickname': 'Bob',
      'friend_avatar_url': 'https://cdn.example.com/bob.jpg',
      'friend_status_message': 'Hello!',
      'status': 'accepted',
    };

    group('fromJson', () {
      test('parses all fields from valid json', () {
        final model = FriendModel.fromJson(validJson);

        expect(model.id, 'friendship-1');
        expect(model.friendId, 'user-bob');
        expect(model.friendNickname, 'Bob');
        expect(model.friendAvatarUrl, 'https://cdn.example.com/bob.jpg');
        expect(model.friendStatusMessage, 'Hello!');
        expect(model.status, FriendStatus.accepted);
      });

      test('parses status "pending" as FriendStatus.pending', () {
        final json = {...validJson, 'status': 'pending'};

        final model = FriendModel.fromJson(json);

        expect(model.status, FriendStatus.pending);
      });

      test('parses status "rejected" as FriendStatus.rejected', () {
        final json = {...validJson, 'status': 'rejected'};

        final model = FriendModel.fromJson(json);

        expect(model.status, FriendStatus.rejected);
      });

      test('defaults to FriendStatus.accepted when status is null', () {
        final json = {...validJson, 'status': null};

        final model = FriendModel.fromJson(json);

        expect(model.status, FriendStatus.accepted);
      });

      test('defaults to FriendStatus.accepted when status is unknown string', () {
        final json = {...validJson, 'status': 'blocked'};

        final model = FriendModel.fromJson(json);

        expect(model.status, FriendStatus.accepted);
      });

      test('sets friendAvatarUrl to null when key value is null', () {
        final json = {...validJson, 'friend_avatar_url': null};

        final model = FriendModel.fromJson(json);

        expect(model.friendAvatarUrl, isNull);
      });

      test('sets friendStatusMessage to null when key value is null', () {
        final json = {...validJson, 'friend_status_message': null};

        final model = FriendModel.fromJson(json);

        expect(model.friendStatusMessage, isNull);
      });
    });

    group('toJson', () {
      test('serializes all fields correctly', () {
        const model = FriendModel(
          id: 'friendship-1',
          friendId: 'user-bob',
          friendNickname: 'Bob',
          friendAvatarUrl: 'https://cdn.example.com/bob.jpg',
          friendStatusMessage: 'Hello!',
          status: FriendStatus.accepted,
        );

        final json = model.toJson();

        expect(json['id'], 'friendship-1');
        expect(json['friend_id'], 'user-bob');
        expect(json['friend_nickname'], 'Bob');
        expect(json['friend_avatar_url'], 'https://cdn.example.com/bob.jpg');
        expect(json['friend_status_message'], 'Hello!');
        expect(json['status'], 'accepted');
      });

      test('serializes FriendStatus.pending as string "pending"', () {
        const model = FriendModel(
          id: 'f2',
          friendId: 'user-carol',
          friendNickname: 'Carol',
          status: FriendStatus.pending,
        );

        expect(model.toJson()['status'], 'pending');
      });

      test('serializes FriendStatus.rejected as string "rejected"', () {
        const model = FriendModel(
          id: 'f3',
          friendId: 'user-dave',
          friendNickname: 'Dave',
          status: FriendStatus.rejected,
        );

        expect(model.toJson()['status'], 'rejected');
      });

      test('serializes null optional fields as null', () {
        const model = FriendModel(
          id: 'f4',
          friendId: 'user-eve',
          friendNickname: 'Eve',
        );

        final json = model.toJson();

        expect(json['friend_avatar_url'], isNull);
        expect(json['friend_status_message'], isNull);
      });

      test('round-trips through fromJson then toJson without data loss', () {
        final original = {
          'id': 'f5',
          'friend_id': 'user-frank',
          'friend_nickname': 'Frank',
          'friend_avatar_url': null,
          'friend_status_message': null,
          'status': 'accepted',
        };

        final json = FriendModel.fromJson(original).toJson();

        expect(json, equals(original));
      });
    });
  });
}
