import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import '../../../auth/presentation/providers/auth_provider.dart';
import '../../../../core/services/user_movie_service.dart';

// Provider para buscar filmes assistidos
final watchedMoviesProvider = FutureProvider<List<dynamic>>((ref) async {
  return await UserMovieService.getWatchedMovies();
});

class WatchedMoviesPage extends ConsumerWidget {
  const WatchedMoviesPage({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final watchedAsync = ref.watch(watchedMoviesProvider);

    return Scaffold(
      appBar: AppBar(
        leading: IconButton(
          icon: const Icon(Icons.arrow_back),
          onPressed: () => context.go('/home'),
        ),
        title: const Text('Assistidos'),
        actions: [
          IconButton(
            icon: const Icon(Icons.refresh),
            onPressed: () => ref.refresh(watchedMoviesProvider),
          ),
          Builder(
            builder: (context) => IconButton(
              icon: const Icon(Icons.menu),
              onPressed: () => Scaffold.of(context).openEndDrawer(),
            ),
          ),
        ],
      ),
      endDrawer: _buildDrawer(context, ref),
      body: watchedAsync.when(
        data: (watched) {
          if (watched.isEmpty) {
            return const Center(
              child: Column(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  Icon(
                    Icons.check_circle_outline,
                    size: 64,
                    color: Colors.grey,
                  ),
                  SizedBox(height: 16),
                  Text(
                    'Nenhum filme assistido ainda',
                    style: TextStyle(fontSize: 18, color: Colors.grey),
                  ),
                ],
              ),
            );
          }

          return GridView.builder(
            padding: const EdgeInsets.all(16),
            gridDelegate: const SliverGridDelegateWithFixedCrossAxisCount(
              crossAxisCount: 2,
              childAspectRatio: 0.7,
              crossAxisSpacing: 16,
              mainAxisSpacing: 16,
            ),
            itemCount: watched.length,
            itemBuilder: (context, index) {
              final movie = watched[index];
              return _buildMovieCard(context, movie);
            },
          );
        },
        loading: () => const Center(child: CircularProgressIndicator()),
        error: (error, stack) => Center(
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              const Icon(Icons.error_outline, size: 48, color: Colors.red),
              const SizedBox(height: 16),
              Text('Erro ao carregar assistidos: $error'),
              const SizedBox(height: 16),
              ElevatedButton(
                onPressed: () => ref.refresh(watchedMoviesProvider),
                child: const Text('Tentar Novamente'),
              ),
            ],
          ),
        ),
      ),
    );
  }

  Widget _buildMovieCard(BuildContext context, dynamic data) {
    print('🎬 Building watched movie card: $data');

    // Os dados do filme estão dentro do objeto 'movie'
    final movie = data['movie'] ?? {};

    final movieId = movie['id'] ?? '';
    final externalApiId = movie['external_api_id'] ?? '';
    final title = movie['title'] ?? 'Sem título';
    final posterUrl = movie['poster_url'] ?? '';

    print('🎬 Movie ID: $movieId, External: $externalApiId');
    print('🎬 Title: $title');
    print('🎬 Poster URL: $posterUrl');

    return GestureDetector(
      onTap: () {
        if (movieId.isNotEmpty) {
          context.go('/movie/$movieId?externalApiId=$externalApiId');
        }
      },
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        mainAxisSize: MainAxisSize.min,
        children: [
          ClipRRect(
            borderRadius: BorderRadius.circular(8),
            child: posterUrl.isNotEmpty
                ? Image.network(
                    posterUrl,
                    height: 195,
                    width: double.infinity,
                    fit: BoxFit.cover,
                    errorBuilder: (context, error, stackTrace) {
                      return _buildPlaceholder();
                    },
                  )
                : _buildPlaceholder(),
          ),
          const SizedBox(height: 4),
          Text(
            title,
            maxLines: 1,
            overflow: TextOverflow.ellipsis,
            style: const TextStyle(fontWeight: FontWeight.w600, fontSize: 12),
          ),
        ],
      ),
    );
  }

  Widget _buildPlaceholder() {
    return Container(
      height: 195,
      width: double.infinity,
      color: Colors.grey[800],
      child: const Center(
        child: Icon(Icons.movie, size: 50, color: Colors.grey),
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
            selected: true,
            onTap: () => Navigator.pop(context),
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
}
