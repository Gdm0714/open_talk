import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:mocktail/mocktail.dart';
import 'package:dio/dio.dart';

import 'package:open_talk/core/network/api_client.dart';
import 'package:open_talk/features/home/data/models/friend_model.dart';
import 'package:open_talk/features/home/domain/providers/home_provider.dart';
import 'package:open_talk/features/home/presentation/screens/add_friend_screen.dart';

class MockDio extends Mock implements Dio {}

class MockFriendListNotifier extends FriendListNotifier {
  @override
  Future<List<FriendModel>> build() async => const [];

  @override
  Future<void> sendRequest(String userId) async {}
}

void main() {
  late MockDio mockDio;

  setUp(() {
    mockDio = MockDio();
  });

  Widget buildSubject() {
    return ProviderScope(
      overrides: [
        apiClientProvider.overrideWithValue(mockDio),
        friendListProvider.overrideWith(MockFriendListNotifier.new),
        // searchUsersProvider returns empty by default (query is empty on start)
      ],
      child: const MaterialApp(
        home: AddFriendScreen(),
      ),
    );
  }

  testWidgets('renders search bar with hint text', (tester) async {
    await tester.pumpWidget(buildSubject());
    await tester.pumpAndSettle();

    expect(find.byType(TextField), findsOneWidget);
    expect(find.text('이메일 또는 닉네임으로 검색'), findsOneWidget);
  });

  testWidgets('renders search icon in search bar', (tester) async {
    await tester.pumpWidget(buildSubject());
    await tester.pumpAndSettle();

    expect(find.byIcon(Icons.search), findsOneWidget);
  });

  testWidgets('shows empty state prompt when query is blank', (tester) async {
    await tester.pumpWidget(buildSubject());
    await tester.pumpAndSettle();

    expect(find.text('이메일 또는 닉네임을 입력하세요'), findsOneWidget);
  });

  testWidgets('shows app bar title 친구 추가', (tester) async {
    await tester.pumpWidget(buildSubject());
    await tester.pumpAndSettle();

    expect(find.text('친구 추가'), findsOneWidget);
  });

  testWidgets('typing into search field updates query', (tester) async {
    await tester.pumpWidget(buildSubject());
    await tester.pumpAndSettle();

    final searchField = find.byType(TextField);
    await tester.enterText(searchField, 'alice');
    await tester.pump();

    // Empty state message should disappear once query is non-empty
    expect(find.text('이메일 또는 닉네임을 입력하세요'), findsNothing);
  });

  testWidgets('clear button appears when query is non-empty', (tester) async {
    await tester.pumpWidget(buildSubject());
    await tester.pumpAndSettle();

    final searchField = find.byType(TextField);
    await tester.enterText(searchField, 'alice');
    await tester.pump();

    expect(find.byIcon(Icons.clear), findsOneWidget);
  });

  testWidgets('tapping clear button resets query', (tester) async {
    await tester.pumpWidget(buildSubject());
    await tester.pumpAndSettle();

    final searchField = find.byType(TextField);
    await tester.enterText(searchField, 'alice');
    await tester.pump();

    await tester.tap(find.byIcon(Icons.clear));
    await tester.pump();

    expect(find.text('이메일 또는 닉네임을 입력하세요'), findsOneWidget);
  });
}
