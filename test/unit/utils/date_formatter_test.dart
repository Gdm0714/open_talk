import 'package:flutter_test/flutter_test.dart';
import 'package:intl/date_symbol_data_local.dart';

import 'package:open_talk/core/utils/date_formatter.dart';

void main() {
  setUpAll(() async {
    // Initialize locale data required by intl/DateFormat with 'ko_KR'
    await initializeDateFormatting('ko_KR', null);
  });

  group('DateFormatter', () {
    group('formatMessageTime', () {
      test('returns time string (오전/오후) for a message sent today', () {
        final now = DateTime.now();
        // Pick a time that is clearly today
        final todayMessage = DateTime(now.year, now.month, now.day, 9, 5);

        final result = DateFormatter.formatMessageTime(todayMessage);

        // Should NOT be '어제' or a date; should contain a digit (the time)
        expect(result, isNot('어제'));
        expect(result.contains(RegExp(r'\d')), isTrue);
      });

      test('returns "어제" for a message sent yesterday', () {
        final yesterday = DateTime.now().subtract(const Duration(days: 1));

        final result = DateFormatter.formatMessageTime(yesterday);

        expect(result, '어제');
      });

      test('returns day-of-week string for a message sent 2-6 days ago', () {
        final twoDaysAgo = DateTime.now().subtract(const Duration(days: 2));

        final result = DateFormatter.formatMessageTime(twoDaysAgo);

        // Should be a Korean weekday name, not a date format
        expect(result, isNot('어제'));
        expect(result.contains(RegExp(r'\d')), isFalse);
      });

      test('returns "M월 d일" format for same-year message older than 7 days', () {
        final now = DateTime.now();
        // Guarantee we are in the same year but more than 7 days ago
        final oldDate = DateTime(now.year, 1, 1);
        // Only run this sub-test when the old date is genuinely > 7 days away
        if (now.difference(oldDate).inDays > 7) {
          final result = DateFormatter.formatMessageTime(oldDate);
          expect(result, contains('월'));
          expect(result, contains('일'));
        }
      });

      test('returns "yyyy.M.d" format for a message from a previous year', () {
        final lastYear = DateTime(DateTime.now().year - 1, 6, 15);

        final result = DateFormatter.formatMessageTime(lastYear);

        // e.g. "2023.6.15"
        expect(result, matches(RegExp(r'^\d{4}\.\d+\.\d+$')));
      });
    });

    group('formatChatListTime', () {
      test('returns time string for a message sent today', () {
        final now = DateTime.now();
        final todayMessage = DateTime(now.year, now.month, now.day, 14, 30);

        final result = DateFormatter.formatChatListTime(todayMessage);

        expect(result.contains(RegExp(r'\d')), isTrue);
        expect(result, isNot('어제'));
      });

      test('returns "어제" for a message sent yesterday', () {
        final yesterday = DateTime.now().subtract(const Duration(days: 1));

        final result = DateFormatter.formatChatListTime(yesterday);

        expect(result, '어제');
      });

      test('returns "M/d" format for same-year message older than 1 day', () {
        final now = DateTime.now();
        final oldDate = DateTime(now.year, 1, 5);
        if (now.difference(oldDate).inDays > 1) {
          final result = DateFormatter.formatChatListTime(oldDate);
          // e.g. "1/5"
          expect(result, matches(RegExp(r'^\d+/\d+$')));
        }
      });

      test('returns "yy/M/d" format for a message from a previous year', () {
        final lastYear = DateTime(DateTime.now().year - 1, 3, 20);

        final result = DateFormatter.formatChatListTime(lastYear);

        // e.g. "23/3/20"
        expect(result, matches(RegExp(r'^\d{2}/\d+/\d+$')));
      });
    });
  });
}
