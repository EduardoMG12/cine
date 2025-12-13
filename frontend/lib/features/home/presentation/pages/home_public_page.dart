import 'dart:async';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import '../providers/movie_provider.dart';
import '../../../movies/data/models/movie_model.dart';

class HomePublicPage extends ConsumerStatefulWidget {
  const HomePublicPage({super.key});

  @override
  ConsumerState<HomePublicPage> createState() => _HomePublicPageState();
}

class _HomePublicPageState extends ConsumerState<HomePublicPage> {
  final _searchController = TextEditingController();
  Timer? _debounce;

  @override
  void dispose() {
    _searchController.dispose();
    _debounce?.cancel();
    super.dispose();
  }

  void _onSearchChanged(String query) {
    if (_debounce?.isActive ?? false) _debounce!.cancel();

    _debounce = Timer(const Duration(milliseconds: 500), () {
      ref.read(searchQueryProvider.notifier).setQuery(query);
      ref.read(searchPageProvider.notifier).reset();
    });
  }

  @override
  Widget build(BuildContext context) {
    final moviesAsync = ref.watch(trendingMoviesProvider);
    final searchQuery = ref.watch(searchQueryProvider);
    final currentPage = ref.watch(searchPageProvider);

    // Só busca se tiver query
    final searchResults = searchQuery.isNotEmpty
        ? ref.watch(
            searchMoviesProvider((query: searchQuery, page: currentPage)),
          )
        : const AsyncValue<List<MovieModel>>.data([]);

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
      body: Column(
        children: [
          // Barra de pesquisa
          Padding(
            padding: const EdgeInsets.all(16),
            child: TextField(
              controller: _searchController,
              onChanged: _onSearchChanged,
              decoration: InputDecoration(
                hintText: 'Pesquisar filmes...',
                prefixIcon: const Icon(Icons.search),
                suffixIcon: searchQuery.isNotEmpty
                    ? IconButton(
                        icon: const Icon(Icons.clear),
                        onPressed: () {
                          _searchController.clear();
                          ref.read(searchQueryProvider.notifier).setQuery('');
                          ref.read(searchPageProvider.notifier).reset();
                        },
                      )
                    : null,
                border: OutlineInputBorder(
                  borderRadius: BorderRadius.circular(12),
                ),
                filled: true,
              ),
            ),
          ),

          // Resultados
          Expanded(
            child: searchQuery.isNotEmpty
                ? _buildSearchResults(context, searchResults)
                : moviesAsync.when(
                    data: (movies) => _buildMovieCarousels(context, movies),
                    loading: () =>
                        const Center(child: CircularProgressIndicator()),
                    error: (error, stack) => Center(
                      child: Column(
                        mainAxisAlignment: MainAxisAlignment.center,
                        children: [
                          const Icon(
                            Icons.error_outline,
                            size: 48,
                            color: Colors.red,
                          ),
                          const SizedBox(height: 16),
                          Text('Erro ao carregar filmes: $error'),
                          const SizedBox(height: 16),
                          ElevatedButton(
                            onPressed: () =>
                                ref.refresh(trendingMoviesProvider),
                            child: const Text('Tentar Novamente'),
                          ),
                        ],
                      ),
                    ),
                  ),
          ),

          // Paginação (só aparece durante pesquisa)
          if (searchQuery.isNotEmpty)
            Container(
              padding: const EdgeInsets.all(16),
              child: Row(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  IconButton(
                    icon: const Icon(Icons.arrow_back),
                    onPressed: currentPage > 1
                        ? () {
                            ref
                                .read(searchPageProvider.notifier)
                                .previousPage();
                          }
                        : null,
                  ),
                  const SizedBox(width: 16),
                  Text(
                    'Página $currentPage',
                    style: const TextStyle(
                      fontSize: 16,
                      fontWeight: FontWeight.w600,
                    ),
                  ),
                  const SizedBox(width: 16),
                  IconButton(
                    icon: const Icon(Icons.arrow_forward),
                    onPressed: () {
                      ref.read(searchPageProvider.notifier).nextPage();
                    },
                  ),
                ],
              ),
            ),
        ],
      ),
    );
  }

  Widget _buildSearchResults(
    BuildContext context,
    AsyncValue<List<MovieModel>> searchResults,
  ) {
    return searchResults.when(
      data: (movies) {
        if (movies.isEmpty) {
          return const Center(
            child: Column(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                Icon(Icons.search_off, size: 64, color: Colors.grey),
                SizedBox(height: 16),
                Text(
                  'Nenhum filme encontrado',
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
            childAspectRatio: 0.6,
            crossAxisSpacing: 12,
            mainAxisSpacing: 12,
          ),
          itemCount: movies.length,
          itemBuilder: (context, index) {
            final movie = movies[index];
            return _buildSearchMovieCard(context, movie);
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
            Text('Erro na pesquisa: $error'),
          ],
        ),
      ),
    );
  }

  Widget _buildSearchMovieCard(BuildContext context, MovieModel movie) {
    return GestureDetector(
      onTap: () => context.go('/movie/${movie.externalApiId}'),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Expanded(
            child: ClipRRect(
              borderRadius: BorderRadius.circular(8),
              child: movie.posterUrl != null && movie.posterUrl!.isNotEmpty
                  ? Image.network(
                      movie.posterUrl!,
                      width: double.infinity,
                      fit: BoxFit.cover,
                      errorBuilder: (context, error, stackTrace) {
                        return _buildPlaceholder();
                      },
                    )
                  : _buildPlaceholder(),
            ),
          ),
          const SizedBox(height: 8),
          Text(
            movie.title,
            maxLines: 2,
            overflow: TextOverflow.ellipsis,
            style: const TextStyle(fontWeight: FontWeight.w600, fontSize: 13),
          ),
          if (movie.overview.isNotEmpty)
            Text(
              movie.overview,
              maxLines: 1,
              overflow: TextOverflow.ellipsis,
              style: TextStyle(fontSize: 11, color: Colors.grey[400]),
            ),
        ],
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
          height: 240,
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
        margin: const EdgeInsets.symmetric(horizontal: 4, vertical: 2),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          mainAxisSize: MainAxisSize.min,
          children: [
            // Poster do filme
            ClipRRect(
              borderRadius: BorderRadius.circular(8),
              child: movie.posterUrl != null && movie.posterUrl!.isNotEmpty
                  ? Image.network(
                      movie.posterUrl!,
                      height: 195,
                      width: 140,
                      fit: BoxFit.cover,
                      errorBuilder: (context, error, stackTrace) {
                        return _buildPlaceholder();
                      },
                    )
                  : _buildPlaceholder(),
            ),
            const SizedBox(height: 4),
            // Título do filme
            Text(
              movie.title,
              maxLines: 1,
              overflow: TextOverflow.ellipsis,
              style: const TextStyle(fontWeight: FontWeight.w600, fontSize: 12),
            ),
            // Ano
            if (movie.overview.isNotEmpty)
              Text(
                movie.overview,
                maxLines: 1,
                overflow: TextOverflow.ellipsis,
                style: TextStyle(fontSize: 10, color: Colors.grey[400]),
              ),
          ],
        ),
      ),
    );
  }

  Widget _buildPlaceholder() {
    return Container(
      height: 195,
      width: 140,
      decoration: BoxDecoration(
        color: Colors.grey[800],
        borderRadius: BorderRadius.circular(8),
      ),
      child: const Icon(Icons.movie, size: 48, color: Colors.grey),
    );
  }
}
