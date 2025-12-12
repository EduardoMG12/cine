import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

class HomePublicPage extends StatelessWidget {
  const HomePublicPage({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('CineVerse'),
        actions: [
          TextButton(
            onPressed: () => context.go('/auth/login'),
            child: const Text('Login'),
          ),
        ],
      ),
      body: Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            const Icon(Icons.movie_filter, size: 100, color: Color(0xFFE50914)),
            const SizedBox(height: 24),
            Text('Welcome to CineVerse', style: Theme.of(context).textTheme.displayMedium),
            const SizedBox(height: 16),
            const Text('Discover, share, and match movies with friends'),
            const SizedBox(height: 32),
            ElevatedButton(
              onPressed: () => context.go('/auth/login'),
              child: const Text('Get Started'),
            ),
          ],
        ),
      ),
    );
  }
}
