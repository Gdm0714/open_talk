import 'package:flutter_test/flutter_test.dart';

import 'package:open_talk/core/utils/validators.dart';

void main() {
  group('Validators', () {
    group('validateEmail', () {
      test('returns null for a valid email address', () {
        expect(Validators.validateEmail('alice@example.com'), isNull);
      });

      test('returns null for email with subdomain', () {
        expect(Validators.validateEmail('user@mail.example.co'), isNull);
      });

      test('returns error message when value is null', () {
        expect(Validators.validateEmail(null), isNotNull);
      });

      test('returns error message when value is empty string', () {
        expect(Validators.validateEmail(''), isNotNull);
      });

      test('returns error message when value is whitespace only', () {
        expect(Validators.validateEmail('   '), isNotNull);
      });

      test('returns error message when @ symbol is missing', () {
        expect(Validators.validateEmail('notanemail.com'), isNotNull);
      });

      test('returns error message when domain is missing after @', () {
        expect(Validators.validateEmail('user@'), isNotNull);
      });

      test('returns error message for email with spaces', () {
        expect(Validators.validateEmail('user name@example.com'), isNotNull);
      });

      test('returns error message when TLD is missing', () {
        expect(Validators.validateEmail('user@domain'), isNotNull);
      });

      test('empty message text starts with expected Korean prompt', () {
        final result = Validators.validateEmail('');
        expect(result, contains('이메일'));
      });

      test('invalid format message mentions format', () {
        final result = Validators.validateEmail('bad');
        expect(result, contains('이메일'));
      });
    });

    group('validatePassword', () {
      test('returns null for password with 8+ chars including letter and digit', () {
        expect(Validators.validatePassword('abcd1234'), isNull);
      });

      test('returns null for longer valid password with mixed chars', () {
        expect(Validators.validatePassword('MySecure1Pass'), isNull);
      });

      test('returns error message when value is null', () {
        expect(Validators.validatePassword(null), isNotNull);
      });

      test('returns error message when value is empty string', () {
        expect(Validators.validatePassword(''), isNotNull);
      });

      test('returns error message when password is shorter than 8 characters', () {
        expect(Validators.validatePassword('abc123'), isNotNull);
      });

      test('returns error message when password is exactly 7 characters', () {
        expect(Validators.validatePassword('abcd123'), isNotNull);
      });

      test('returns error message when password contains no letters', () {
        expect(Validators.validatePassword('12345678'), isNotNull);
      });

      test('returns error message when password contains no digits', () {
        expect(Validators.validatePassword('abcdefgh'), isNotNull);
      });

      test('returns null for exactly 8 chars with letter and digit', () {
        expect(Validators.validatePassword('aaaaa1aa'), isNull);
      });
    });

    group('validateNickname', () {
      test('returns null for a valid 2-character nickname', () {
        expect(Validators.validateNickname('AB'), isNull);
      });

      test('returns null for a valid nickname within length bounds', () {
        expect(Validators.validateNickname('OpenTalkUser'), isNull);
      });

      test('returns null for a nickname with exactly 20 characters', () {
        expect(Validators.validateNickname('A' * 20), isNull);
      });

      test('returns error message when value is null', () {
        expect(Validators.validateNickname(null), isNotNull);
      });

      test('returns error message when value is empty string', () {
        expect(Validators.validateNickname(''), isNotNull);
      });

      test('returns error message when value is whitespace only', () {
        expect(Validators.validateNickname('  '), isNotNull);
      });

      test('returns error message when nickname is 1 character', () {
        expect(Validators.validateNickname('A'), isNotNull);
      });

      test('returns error message when nickname exceeds 20 characters', () {
        expect(Validators.validateNickname('A' * 21), isNotNull);
      });

      test('trims whitespace before checking length for too-short case', () {
        // ' A ' trims to 'A' which is 1 char — should fail
        expect(Validators.validateNickname(' A '), isNotNull);
      });

      test('returns null for nickname with exactly 2 non-whitespace chars', () {
        // '  AB  ' trims to 'AB' (2 chars) — should pass
        expect(Validators.validateNickname('  AB  '), isNull);
      });
    });

    group('validateRequired', () {
      test('returns null when value is non-empty', () {
        expect(Validators.validateRequired('something', '필드'), isNull);
      });

      test('returns error message when value is null', () {
        expect(Validators.validateRequired(null, '이름'), isNotNull);
      });

      test('returns error message when value is empty', () {
        expect(Validators.validateRequired('', '이름'), isNotNull);
      });

      test('error message contains the field name', () {
        final result = Validators.validateRequired('', '닉네임');
        expect(result, contains('닉네임'));
      });
    });
  });
}
