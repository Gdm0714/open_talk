import 'package:intl/intl.dart';

class DateFormatter {
  DateFormatter._();

  static String formatMessageTime(DateTime dateTime) {
    final now = DateTime.now();
    final today = DateTime(now.year, now.month, now.day);
    final messageDate = DateTime(dateTime.year, dateTime.month, dateTime.day);
    final difference = today.difference(messageDate).inDays;

    if (difference == 0) {
      return DateFormat('a h:mm', 'ko_KR').format(dateTime);
    } else if (difference == 1) {
      return '어제';
    } else if (difference < 7) {
      return DateFormat('EEEE', 'ko_KR').format(dateTime);
    } else if (dateTime.year == now.year) {
      return DateFormat('M월 d일').format(dateTime);
    } else {
      return DateFormat('yyyy.M.d').format(dateTime);
    }
  }

  static String formatChatListTime(DateTime dateTime) {
    final now = DateTime.now();
    final today = DateTime(now.year, now.month, now.day);
    final messageDate = DateTime(dateTime.year, dateTime.month, dateTime.day);
    final difference = today.difference(messageDate).inDays;

    if (difference == 0) {
      return DateFormat('a h:mm', 'ko_KR').format(dateTime);
    } else if (difference == 1) {
      return '어제';
    } else if (dateTime.year == now.year) {
      return DateFormat('M/d').format(dateTime);
    } else {
      return DateFormat('yy/M/d').format(dateTime);
    }
  }

  static String formatFullDate(DateTime dateTime) {
    return DateFormat('yyyy년 M월 d일 a h:mm', 'ko_KR').format(dateTime);
  }
}
