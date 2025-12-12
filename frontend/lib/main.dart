import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'core/router/app_router.dart';
import 'core/theme/app_theme.dart';
import 'core/services/api_service.dart';

void main() {
  // Initialize API service
  ApiService.init();
  
  runApp(const ProviderScope(child: CineVerseApp()));
}

class CineVerseApp extends ConsumerWidget {
  const CineVerseApp({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final router = ref.watch(routerProvider);

    return MaterialApp.router(
      title: 'CineVerse',
      theme: AppTheme.darkTheme,
      debugShowCheckedModeBanner: false,
      routerConfig: router,
    );
  }
}
