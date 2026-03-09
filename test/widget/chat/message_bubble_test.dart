import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:intl/date_symbol_data_local.dart';

import 'package:open_talk/features/chat/data/models/message_model.dart';
import 'package:open_talk/features/chat/presentation/widgets/message_bubble.dart';

import '../../helpers/test_helpers.dart';

MessageModel _makeMessage({
  String id = 'msg-1',
  String senderId = 'user-alice',
  String senderNickname = 'Alice',
  String content = 'Hello!',
  MessageType messageType = MessageType.text,
}) {
  return MessageModel(
    id: id,
    chatRoomId: 'room-1',
    senderId: senderId,
    senderNickname: senderNickname,
    content: content,
    messageType: messageType,
    createdAt: DateTime(2024, 6, 1, 10, 30),
  );
}

void main() {
  setUpAll(() async {
    await initializeDateFormatting('ko_KR', null);
  });

  group('MessageBubble', () {
    group('sent message (isMine: true)', () {
      testWidgets('aligns bubble to the right side of the row', (tester) async {
        await tester.pumpWidget(
          buildTestableWidget(
            MessageBubble(
              message: _makeMessage(),
              isMine: true,
            ),
          ),
        );

        final row = tester.widget<Row>(
          find.byWidgetPredicate(
            (w) => w is Row && w.mainAxisAlignment == MainAxisAlignment.end,
          ),
        );
        expect(row.mainAxisAlignment, MainAxisAlignment.end);
      });

      testWidgets('displays the message content text', (tester) async {
        await tester.pumpWidget(
          buildTestableWidget(
            MessageBubble(
              message: _makeMessage(content: 'Hi there!'),
              isMine: true,
            ),
          ),
        );

        expect(find.text('Hi there!'), findsOneWidget);
      });

      testWidgets('does not show sender name for own messages', (tester) async {
        await tester.pumpWidget(
          buildTestableWidget(
            MessageBubble(
              message: _makeMessage(senderNickname: 'Alice'),
              isMine: true,
              showSenderName: true,
            ),
          ),
        );

        // showSenderName only applies to received messages (!isMine)
        expect(find.text('Alice'), findsNothing);
      });
    });

    group('received message (isMine: false)', () {
      testWidgets('aligns bubble to the left side of the row', (tester) async {
        await tester.pumpWidget(
          buildTestableWidget(
            MessageBubble(
              message: _makeMessage(senderId: 'user-bob', senderNickname: 'Bob'),
              isMine: false,
            ),
          ),
        );

        final row = tester.widget<Row>(
          find.byWidgetPredicate(
            (w) => w is Row && w.mainAxisAlignment == MainAxisAlignment.start,
          ),
        );
        expect(row.mainAxisAlignment, MainAxisAlignment.start);
      });

      testWidgets('shows sender name when showSenderName is true', (tester) async {
        await tester.pumpWidget(
          buildTestableWidget(
            MessageBubble(
              message: _makeMessage(senderNickname: 'Bob'),
              isMine: false,
              showSenderName: true,
            ),
          ),
        );

        expect(find.text('Bob'), findsOneWidget);
      });

      testWidgets('does not show sender name when showSenderName is false', (tester) async {
        await tester.pumpWidget(
          buildTestableWidget(
            MessageBubble(
              message: _makeMessage(senderNickname: 'Bob'),
              isMine: false,
              showSenderName: false,
            ),
          ),
        );

        expect(find.text('Bob'), findsNothing);
      });
    });

    group('system message', () {
      testWidgets('renders system message content centered', (tester) async {
        await tester.pumpWidget(
          buildTestableWidget(
            MessageBubble(
              message: _makeMessage(
                content: '채팅방이 생성되었습니다',
                messageType: MessageType.system,
              ),
              isMine: false,
            ),
          ),
        );

        expect(find.text('채팅방이 생성되었습니다'), findsOneWidget);
        // System messages are wrapped in a Center widget
        expect(find.byType(Center), findsOneWidget);
      });
    });
  });
}
