import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:mocktail/mocktail.dart';
import 'package:dio/dio.dart';

import 'package:open_talk/core/theme/app_theme.dart';
import 'package:open_talk/shared/services/auth_storage.dart';

// ---------------------------------------------------------------------------
// Mock classes
// ---------------------------------------------------------------------------

class MockDio extends Mock implements Dio {}

class MockAuthStorage extends Mock implements AuthStorage {}

// ---------------------------------------------------------------------------
// Widget test helper
// ---------------------------------------------------------------------------

/// Wraps [widget] in a MaterialApp > Scaffold > ProviderScope so widget tests
/// have:
///   - A Material ancestor (required by InkWell, TextField, etc.)
///   - Bounded layout constraints (avoids RenderFlex overflow in unit tests)
///   - Optional Riverpod provider [overrides]
Widget buildTestableWidget(
  Widget widget, {
  List<Override> overrides = const [],
}) {
  return ProviderScope(
    overrides: overrides,
    child: MaterialApp(
      theme: AppTheme.light,
      home: Scaffold(
        body: SingleChildScrollView(child: widget),
      ),
    ),
  );
}

/// Helper to pump a widget and wait for all animations/futures to settle.
Future<void> pumpAndSettle(
  WidgetTester tester,
  Widget widget, {
  List<Override> overrides = const [],
}) async {
  await tester.pumpWidget(buildTestableWidget(widget, overrides: overrides));
  await tester.pumpAndSettle();
}
