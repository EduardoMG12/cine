import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import '../providers/movie_provider.dart';
import '../../../movies/data/models/movie_model.dart';

class HomePublicPage extends ConsumerWidget {
  const HomePublicPage({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final moviesAsync = ref.watch(trendingMoviesProvider);

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
      body: moviesAsync.when(
        data: (movies) => _buildMovieCarousels(context, movies),
        loading: () => const Center(child: CircularProgressIndicator()),
        error: (error, stack) => Center(
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              const Icon(Icons.error_outline, size: 48, color: Colors.red),
              const SizedBox(height: 16),
              Text('Erro ao carregar filmes: $error'),
              const SizedBox(height: 16),
              ElevatedButton(
                onPressed: () => ref.refresh(trendingMoviesProvider),
                child: const Text('Tentar Novamente'),
              ),
            ],
          ),
        ),
      ),
    );
  }

  Widget _buildMovieCarousels(
    BuildContext context,
    List<MovieModel> allMovies,
  ) {
    // Divide os filmes em diferentes carroseis
    final recentMovies = allMovies.take(10).toList();
    final popularMovies = allMovies.skip(10).take(10).toList();
    final moreMovies = allMovies.skip(20).take(10).toList();

    return RefreshIndicator(
      onRefresh: () async {
        // Força refresh dos filmes
      },
      child: ListView(
        padding: const EdgeInsets.symmetric(vertical: 16),
        children: [
          // Banner de boas-vindas
          Container(
            margin: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
            padding: const EdgeInsets.all(24),
            decoration: BoxDecoration(
              gradient: const LinearGradient(
                colors: [Color(0xFFE50914), Color(0xFF831010)],
                begin: Alignment.topLeft,
                end: Alignment.bottomRight,
              ),
              borderRadius: BorderRadius.circular(16),
            ),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                const Icon(Icons.movie_filter, size: 48, color: Colors.white),
                const SizedBox(height: 16),
                Text(
                  'Bem-vindo ao CineVerse',
                  style: Theme.of(context).textTheme.headlineSmall?.copyWith(
                    color: Colors.white,
                    fontWeight: FontWeight.bold,
                  ),
                ),
                const SizedBox(height: 8),
                const Text(
                  'Descubra, compartilhe e encontre filmes com amigos',
                  style: TextStyle(color: Colors.white70),
                ),
                const SizedBox(height: 16),
                ElevatedButton(
                  onPressed: () => context.go('/auth/login'),
                  style: ElevatedButton.styleFrom(
                    backgroundColor: Colors.white,
                    foregroundColor: const Color(0xFFE50914),
                  ),
                  child: const Text('Começar'),
                ),
              ],
            ),
          ),

          const SizedBox(height: 24),

          // Carrossel: Filmes Recentes
          if (recentMovies.isNotEmpty)
            _buildMovieCarousel(context, 'Filmes em Alta', recentMovies),

          const SizedBox(height: 24),

          // Carrossel: Filmes Populares
          if (popularMovies.isNotEmpty)
            _buildMovieCarousel(context, 'Populares', popularMovies),

          const SizedBox(height: 24),

          // Carrossel: Mais Filmes
          if (moreMovies.isNotEmpty)
            _buildMovieCarousel(context, 'Descubra Mais', moreMovies),

          const SizedBox(height: 24),
        ],
      ),
    );
  }

  Widget _buildMovieCarousel(
    BuildContext context,
    String title,
    List<MovieModel> movies,
  ) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Padding(
          padding: const EdgeInsets.symmetric(horizontal: 16),
          child: Text(
            title,
            style: Theme.of(
              context,
            ).textTheme.titleLarge?.copyWith(fontWeight: FontWeight.bold),
          ),
        ),
        const SizedBox(height: 12),
        SizedBox(
          height: 220,
          child: ListView.builder(
            scrollDirection: Axis.horizontal,
            padding: const EdgeInsets.symmetric(horizontal: 12),
            itemCount: movies.length,
            itemBuilder: (context, index) {
              final movie = movies[index];
              return _buildMovieCard(context, movie);
            },
          ),
        ),
      ],
    );
  }

  Widget _buildMovieCard(BuildContext context, MovieModel movie) {
    return GestureDetector(
      onTap: () {
        // Navega para a página de detalhes usando o external_api_id
        context.go('/movie/${movie.externalApiId}');
      },
      child: Container(
        width: 140,
        margin: const EdgeInsets.symmetric(horizontal: 4),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            // Poster do filme
            ClipRRect(
              borderRadius: BorderRadius.circular(8),
              child: movie.posterUrl != null && movie.posterUrl!.isNotEmpty
                  ? Image.network(
                      movie.posterUrl!,
                      height: 180,
                      width: 140,
                      fit: BoxFit.cover,
                      errorBuilder: (context, error, stackTrace) {
                        return _buildPlaceholder();
                      },
                    )
                  : _buildPlaceholder(),
            ),
            const SizedBox(height: 8),
            // Título do filme
            Text(
              movie.title,
              maxLines: 1,
              overflow: TextOverflow.ellipsis,
              style: const TextStyle(fontWeight: FontWeight.w600, fontSize: 13),
            ),
            // Ano
            if (movie.overview.isNotEmpty)
              Text(
                movie.overview,
                maxLines: 1,
                overflow: TextOverflow.ellipsis,
                style: TextStyle(fontSize: 11, color: Colors.grey[400]),
              ),
          ],
        ),
      ),
    );
  }

  Widget _buildPlaceholder() {
    return Container(
      height: 180,
      width: 140,
      decoration: BoxDecoration(
        color: Colors.grey[800],
        borderRadius: BorderRadius.circular(8),
      ),
      child: const Icon(Icons.movie, size: 48, color: Colors.grey),
    );
  }
}
