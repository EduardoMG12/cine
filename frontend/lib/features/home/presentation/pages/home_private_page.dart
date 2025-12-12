import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import '../../../auth/presentation/providers/auth_provider.dart';

class HomePrivatePage extends ConsumerWidget {
  const HomePrivatePage({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Home'),
        actions: [
          Builder(
            builder: (context) => IconButton(
              icon: const Icon(Icons.menu),
              onPressed: () => Scaffold.of(context).openEndDrawer(),
            ),
          ),
        ],
      ),
      endDrawer: _buildDrawer(context, ref),
      body: ListView(
        padding: const EdgeInsets.all(16),
        children: [
          _buildQuickAccessCard(
            context,
            'Watch Later',
            Icons.bookmark,
            '/watch-later',
          ),
          _buildQuickAccessCard(
            context,
            'Watched Movies',
            Icons.check_circle,
            '/watched',
          ),
          _buildQuickAccessCard(context, 'Friends', Icons.people, '/friends'),
          _buildQuickAccessCard(
            context,
            'Match Movies',
            Icons.favorite,
            '/match',
          ),
        ],
      ),
    );
  }

  Widget _buildDrawer(BuildContext context, WidgetRef ref) {
    return Drawer(
      child: ListView(
        padding: EdgeInsets.zero,
        children: [
          DrawerHeader(
            decoration: const BoxDecoration(color: Color(0xFFE50914)),
            child: Column(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                const Icon(Icons.movie, size: 64, color: Colors.white),
                const SizedBox(height: 8),
                Text(
                  'Cine',
                  style: Theme.of(context).textTheme.headlineSmall?.copyWith(
                    color: Colors.white,
                    fontWeight: FontWeight.bold,
                  ),
                ),
              ],
            ),
          ),
          ListTile(
            leading: const Icon(Icons.home),
            title: const Text('Home'),
            onTap: () {
              Navigator.pop(context);
              context.go('/home');
            },
          ),
          ListTile(
            leading: const Icon(Icons.bookmark),
            title: const Text('Favoritos'),
            onTap: () {
              Navigator.pop(context);
              context.go('/watch-later');
            },
          ),
          ListTile(
            leading: const Icon(Icons.check_circle),
            title: const Text('Assistidos'),
            onTap: () {
              Navigator.pop(context);
              context.go('/watched');
            },
          ),
          const Divider(),
          ListTile(
            leading: const Icon(Icons.people),
            title: const Text('Amigos'),
            onTap: () {
              Navigator.pop(context);
              context.go('/friends');
            },
          ),
          ListTile(
            leading: const Icon(Icons.favorite),
            title: const Text('Match de Filmes'),
            onTap: () {
              Navigator.pop(context);
              context.go('/match');
            },
          ),
          ListTile(
            leading: const Icon(Icons.person),
            title: const Text('Perfil'),
            onTap: () {
              Navigator.pop(context);
              context.go('/profile');
            },
          ),
          const Divider(),
          ListTile(
            leading: const Icon(Icons.exit_to_app, color: Colors.red),
            title: const Text('Sair', style: TextStyle(color: Colors.red)),
            onTap: () async {
              Navigator.pop(context);
              await ref.read(authStateProvider.notifier).logout();
              if (context.mounted) {
                context.go('/');
              }
            },
          ),
        ],
      ),
    );
  }

  Widget _buildQuickAccessCard(
    BuildContext context,
    String title,
    IconData icon,
    String route,
  ) {
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
