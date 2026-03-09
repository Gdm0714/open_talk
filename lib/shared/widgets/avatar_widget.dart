import 'package:cached_network_image/cached_network_image.dart';
import 'package:flutter/material.dart';

import '../../core/constants/app_colors.dart';

class AvatarWidget extends StatelessWidget {
  final String? imageUrl;
  final String name;
  final double size;
  final bool showOnlineIndicator;
  final bool isOnline;

  const AvatarWidget({
    super.key,
    this.imageUrl,
    required this.name,
    this.size = 48,
    this.showOnlineIndicator = false,
    this.isOnline = false,
  });

  String get _initials {
    if (name.isEmpty) return '?';
    final parts = name.trim().split(' ');
    if (parts.length >= 2) {
      return '${parts[0][0]}${parts[1][0]}'.toUpperCase();
    }
    return name[0].toUpperCase();
  }

  @override
  Widget build(BuildContext context) {
    return Stack(
      children: [
        CircleAvatar(
          radius: size / 2,
          backgroundColor: AppColors.primaryLight.withValues(alpha: 0.3),
          backgroundImage:
              imageUrl != null && imageUrl!.isNotEmpty
                  ? CachedNetworkImageProvider(imageUrl!)
                  : null,
          child:
              imageUrl == null || imageUrl!.isEmpty
                  ? Text(
                    _initials,
                    style: TextStyle(
                      color: AppColors.primary,
                      fontSize: size * 0.36,
                      fontWeight: FontWeight.w600,
                    ),
                  )
                  : null,
        ),
        if (showOnlineIndicator)
          Positioned(
            bottom: 0,
            right: 0,
            child: Container(
              width: size * 0.28,
              height: size * 0.28,
              decoration: BoxDecoration(
                color: isOnline ? AppColors.online : AppColors.textHint,
                shape: BoxShape.circle,
                border: Border.all(color: Colors.white, width: 2),
              ),
            ),
          ),
      ],
    );
  }
}
