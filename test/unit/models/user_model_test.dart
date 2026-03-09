import 'package:flutter_test/flutter_test.dart';

import 'package:open_talk/features/auth/data/models/user_model.dart';

void main() {
  group('UserModel', () {
    group('fromJson', () {
      test('parses all required fields from valid json', () {
        final json = {
          'id': 'user-123',
          'email': 'alice@example.com',
          'nickname': 'Alice',
          'avatar_url': 'https://example.com/avatar.png',
          'status_message': 'Hello!',
        };

        final model = UserModel.fromJson(json);

        expect(model.id, 'user-123');
        expect(model.email, 'alice@example.com');
        expect(model.nickname, 'Alice');
        expect(model.avatarUrl, 'https://example.com/avatar.png');
        expect(model.statusMessage, 'Hello!');
      });

      test('sets avatarUrl to null when key is absent', () {
        final json = {
          'id': 'user-456',
          'email': 'bob@example.com',
          'nickname': 'Bob',
          'avatar_url': null,
          'status_message': null,
        };

        final model = UserModel.fromJson(json);

        expect(model.avatarUrl, isNull);
        expect(model.statusMessage, isNull);
      });

      test('sets statusMessage to null when key is absent', () {
        final json = {
          'id': 'user-789',
          'email': 'carol@example.com',
          'nickname': 'Carol',
        };

        final model = UserModel.fromJson(json);

        expect(model.statusMessage, isNull);
      });
    });

    group('toJson', () {
      test('produces correct map with all fields populated', () {
        const model = UserModel(
          id: 'user-123',
          email: 'alice@example.com',
          nickname: 'Alice',
          avatarUrl: 'https://example.com/avatar.png',
          statusMessage: 'Hello!',
        );

        final json = model.toJson();

        expect(json['id'], 'user-123');
        expect(json['email'], 'alice@example.com');
        expect(json['nickname'], 'Alice');
        expect(json['avatar_url'], 'https://example.com/avatar.png');
        expect(json['status_message'], 'Hello!');
      });

      test('produces correct map when optional fields are null', () {
        const model = UserModel(
          id: 'user-456',
          email: 'bob@example.com',
          nickname: 'Bob',
        );

        final json = model.toJson();

        expect(json['avatar_url'], isNull);
        expect(json['status_message'], isNull);
      });

      test('round-trips through fromJson then toJson without data loss', () {
        final original = {
          'id': 'user-999',
          'email': 'test@test.com',
          'nickname': 'Tester',
          'avatar_url': 'https://cdn.example.com/img.jpg',
          'status_message': 'Working hard',
        };

        final json = UserModel.fromJson(original).toJson();

        expect(json, equals(original));
      });
    });

    group('copyWith', () {
      test('returns new instance with updated nickname only', () {
        const original = UserModel(
          id: 'u1',
          email: 'a@b.com',
          nickname: 'Old',
        );

        final updated = original.copyWith(nickname: 'New');

        expect(updated.nickname, 'New');
        expect(updated.id, original.id);
        expect(updated.email, original.email);
      });
    });
  });
}
