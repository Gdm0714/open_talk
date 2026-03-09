import 'package:flutter/material.dart';

import '../../../../core/constants/app_colors.dart';
import '../../../../core/utils/date_formatter.dart';
import '../../../../shared/widgets/avatar_widget.dart';
import '../../data/models/chat_room_model.dart';

class ChatRoomTile extends StatelessWidget {
  final ChatRoomModel chatRoom;
  final String currentUserId;
  final VoidCallback? onTap;

  const ChatRoomTile({
    super.key,
    required this.chatRoom,
    required this.currentUserId,
    this.onTap,
  });

  @override
  Widget build(BuildContext context) {
    final displayName = chatRoom.displayName(currentUserId);
    final otherMember = chatRoom.type == ChatRoomType.direct
        ? chatRoom.members.where((m) => m.id != currentUserId).firstOrNull
        : null;

    return InkWell(
      onTap: onTap,
      child: Padding(
        padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 12),
        child: Row(
          children: [
            // Avatar
            AvatarWidget(
              name: displayName,
              imageUrl: otherMember?.avatarUrl,
              size: 52,
            ),
            const SizedBox(width: 14),

            // Content
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  // Name and time
                  Row(
                    children: [
                      Expanded(
                        child: Text(
                          displayName,
                          style: Theme.of(context).textTheme.titleMedium,
                          maxLines: 1,
                          overflow: TextOverflow.ellipsis,
                        ),
                      ),
                      if (chatRoom.lastMessageAt != null)
                        Text(
                          DateFormatter.formatChatListTime(
                            chatRoom.lastMessageAt!,
                          ),
                          style:
                              Theme.of(context).textTheme.bodySmall?.copyWith(
                                    color: AppColors.textHint,
                                    fontSize: 12,
                                  ),
                        ),
                    ],
                  ),
                  const SizedBox(height: 4),

                  // Last message and unread
                  Row(
                    children: [
                      Expanded(
                        child: Text(
                          chatRoom.lastMessage ?? '',
                          style:
                              Theme.of(context).textTheme.bodyMedium?.copyWith(
                                    color: AppColors.textSecondary,
                                  ),
                          maxLines: 1,
                          overflow: TextOverflow.ellipsis,
                        ),
                      ),
                      if (chatRoom.unreadCount > 0) ...[
                        const SizedBox(width: 8),
                        Container(
                          padding: const EdgeInsets.symmetric(
                            horizontal: 8,
                            vertical: 2,
                          ),
                          decoration: BoxDecoration(
                            color: AppColors.unreadBadge,
                            borderRadius: BorderRadius.circular(10),
                          ),
                          child: Text(
                            chatRoom.unreadCount > 99
                                ? '99+'
                                : '${chatRoom.unreadCount}',
                            style: const TextStyle(
                              color: Colors.white,
                              fontSize: 12,
                              fontWeight: FontWeight.w600,
                            ),
                          ),
                        ),
                      ],
                    ],
                  ),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }
}
