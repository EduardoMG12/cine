import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

class HomePrivatePage extends StatelessWidget {
  const HomePrivatePage({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Home'),
        actions: [
          IconButton(
            icon: const Icon(Icons.person),
            onPressed: () => context.go('/profile'),
          ),
        ],
      ),
      body: ListView(
        padding: const EdgeInsets.all(16),
        children: [
          _buildQuickAccessCard(context, 'Watch Later', Icons.bookmark, '/watch-later'),
          _buildQuickAccessCard(context, 'Watched Movies', Icons.check_circle, '/watched'),
          _buildQuickAccessCard(context, 'Friends', Icons.people, '/friends'),
          _buildQuickAccessCard(context, 'Match Movies', Icons.favorite, '/match'),
        ],
      ),
    );
  }

  Widget _buildQuickAccessCard(BuildContext context, String title, IconData icon, String route) {
    return Card(
      child: ListTile(
        leading: Icon(icon, color: const Color(0xFFE50914)),
        title: Text(title),
        trailing: const Icon(Icons.arrow_forward_ios, size: 16),
        onTap: () => context.go(route),
      ),
    );
  }
}
