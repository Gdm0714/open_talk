import 'package:flutter_test/flutter_test.dart';

import 'package:open_talk/features/auth/data/models/auth_response_model.dart';
import 'package:open_talk/features/auth/data/models/user_model.dart';

void main() {
  group('AuthResponseModel', () {
    final validJson = {
      'token': 'access-abc',
      'user': {
        'id': 'user-1',
        'email': 'alice@example.com',
        'nickname': 'Alice',
        'avatar_url': null,
        'status_message': null,
      },
    };

    group('fromJson', () {
      test('parses accessToken and refreshToken correctly', () {
        final model = AuthResponseModel.fromJson(validJson);

        expect(model.accessToken, 'access-abc');
        expect(model.refreshToken, 'access-abc');
      });

      test('parses nested UserModel correctly', () {
        final model = AuthResponseModel.fromJson(validJson);

        expect(model.user, isA<UserModel>());
        expect(model.user.id, 'user-1');
        expect(model.user.email, 'alice@example.com');
        expect(model.user.nickname, 'Alice');
      });

      test('nested user has null optional fields when json contains null', () {
        final model = AuthResponseModel.fromJson(validJson);

        expect(model.user.avatarUrl, isNull);
        expect(model.user.statusMessage, isNull);
      });

      test('parses nested user with all optional fields populated', () {
        final json = {
          'token': 'tok',
          'user': {
            'id': 'u2',
            'email': 'bob@example.com',
            'nickname': 'Bob',
            'avatar_url': 'https://cdn.example.com/bob.jpg',
            'status_message': 'Hey there',
          },
        };

        final model = AuthResponseModel.fromJson(json);

        expect(model.user.avatarUrl, 'https://cdn.example.com/bob.jpg');
        expect(model.user.statusMessage, 'Hey there');
      });
    });

    group('toJson', () {
      test('produces correct top-level token keys', () {
        const model = AuthResponseModel(
          accessToken: 'access-abc',
          refreshToken: 'refresh-xyz',
          user: UserModel(
            id: 'user-1',
            email: 'alice@example.com',
            nickname: 'Alice',
          ),
        );

        final json = model.toJson();

        expect(json['token'], 'access-abc');
      });

      test('serializes nested user as a map', () {
        const model = AuthResponseModel(
          accessToken: 'tok',
          refreshToken: 'ref',
          user: UserModel(
            id: 'u1',
            email: 'a@b.com',
            nickname: 'A',
          ),
        );

        final json = model.toJson();

        expect(json['user'], isA<Map<String, dynamic>>());
        expect((json['user'] as Map<String, dynamic>)['id'], 'u1');
      });
    });
  });
}
