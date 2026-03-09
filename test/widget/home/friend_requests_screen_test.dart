import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:mocktail/mocktail.dart';
import 'package:dio/dio.dart';

import 'package:open_talk/core/network/api_client.dart';
import 'package:open_talk/features/home/data/models/friend_model.dart';
import 'package:open_talk/features/home/domain/providers/home_provider.dart';
import 'package:open_talk/features/home/presentation/screens/friend_requests_screen.dart';

class MockDio extends Mock implements Dio {}

class MockFriendListNotifier extends FriendListNotifier {
  final List<FriendModel> _friends;
  MockFriendListNotifier(this._friends);

  @override
  Future<List<FriendModel>> build() async => _friends;

  @override
  Future<void> refresh() async {
    state = AsyncValue.data(_friends);
  }

  @override
  Future<void> acceptRequest(String id) async {
    final updated = _friends
        .map((f) => f.id == id
            ? FriendModel(
                id: f.id,
                friendId: f.friendId,
                friendNickname: f.friendNickname,
                friendAvatarUrl: f.friendAvatarUrl,
                friendStatusMessage: f.friendStatusMessage,
                status: FriendStatus.accepted,
              )
            : f)
        .toList();
    state = AsyncValue.data(updated);
  }

  @override
  Future<void> rejectRequest(String id) async {
    final updated = _friends.where((f) => f.id != id).toList();
    state = AsyncValue.data(updated);
  }
}

const _pendingFriend = FriendModel(
  id: 'req-1',
  friendId: 'user-2',
  friendNickname: 'Alice',
  friendStatusMessage: 'Hi there',
  status: FriendStatus.pending,
);

const _acceptedFriend = FriendModel(
  id: 'req-2',
  friendId: 'user-3',
  friendNickname: 'Bob',
  status: FriendStatus.accepted,
);

void main() {
  late MockDio mockDio;

  setUp(() {
    mockDio = MockDio();
  });

  Widget buildSubject(List<FriendModel> friends) {
    return ProviderScope(
      overrides: [
        apiClientProvider.overrideWithValue(mockDio),
        friendListProvider.overrideWith(() => MockFriendListNotifier(friends)),
      ],
      child: const MaterialApp(
        home: FriendRequestsScreen(),
      ),
    );
  }

  testWidgets('renders screen title 친구 요청', (tester) async {
    await tester.pumpWidget(buildSubject([]));
    await tester.pumpAndSettle();

    expect(find.text('친구 요청'), findsOneWidget);
  });

  testWidgets('shows empty state when no pending requests', (tester) async {
    await tester.pumpWidget(buildSubject([_acceptedFriend]));
    await tester.pumpAndSettle();

    expect(find.text('받은 친구 요청이 없습니다'), findsOneWidget);
    expect(find.byIcon(Icons.inbox_outlined), findsOneWidget);
  });

  testWidgets('shows empty state when friend list is empty', (tester) async {
    await tester.pumpWidget(buildSubject([]));
    await tester.pumpAndSettle();

    expect(find.text('받은 친구 요청이 없습니다'), findsOneWidget);
  });

  testWidgets('shows pending request with nickname', (tester) async {
    await tester.pumpWidget(buildSubject([_pendingFriend]));
    await tester.pumpAndSettle();

    expect(find.text('Alice'), findsOneWidget);
  });

  testWidgets('shows accept button for pending request', (tester) async {
    await tester.pumpWidget(buildSubject([_pendingFriend]));
    await tester.pumpAndSettle();

    expect(find.text('수락'), findsOneWidget);
  });

  testWidgets('shows reject button for pending request', (tester) async {
    await tester.pumpWidget(buildSubject([_pendingFriend]));
    await tester.pumpAndSettle();

    expect(find.text('거절'), findsOneWidget);
  });

  testWidgets('does not show accepted friends in list', (tester) async {
    await tester.pumpWidget(buildSubject([_acceptedFriend, _pendingFriend]));
    await tester.pumpAndSettle();

    // Bob is accepted, should not appear as a request tile
    // Alice is pending and should appear
    expect(find.text('Alice'), findsOneWidget);
    expect(find.text('Bob'), findsNothing);
  });

  testWidgets('shows status message when present', (tester) async {
    await tester.pumpWidget(buildSubject([_pendingFriend]));
    await tester.pumpAndSettle();

    expect(find.text('Hi there'), findsOneWidget);
  });
}
