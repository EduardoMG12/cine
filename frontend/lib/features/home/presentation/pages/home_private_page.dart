import 'dart:async';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import '../../../auth/presentation/providers/auth_provider.dart';
import '../providers/movie_provider.dart';
import '../../../movies/data/models/movie_model.dart';

class HomePrivatePage extends ConsumerStatefulWidget {
  const HomePrivatePage({super.key});

  @override
  ConsumerState<HomePrivatePage> createState() => _HomePrivatePageState();
}

class _HomePrivatePageState extends ConsumerState<HomePrivatePage> {
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

    final searchResults = searchQuery.isNotEmpty
        ? ref.watch(
            searchMoviesProvider((query: searchQuery, page: currentPage)),
          )
        : const AsyncValue<List<MovieModel>>.data([]);

    return Scaffold(
      appBar: AppBar(
        title: const Text('CineVerse'),
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

          // Botões de acesso rápido (só aparecem quando não está pesquisando)
          if (searchQuery.isEmpty)
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 16),
              child: Row(
                children: [
                  Expanded(
                    child: _buildQuickAccessButton(
                      context,
                      'Favoritos',
                      Icons.bookmark,
                      '/watch-later',
                    ),
                  ),
                  const SizedBox(width: 12),
                  Expanded(
                    child: _buildQuickAccessButton(
                      context,
                      'Assistidos',
                      Icons.check_circle,
                      '/watched',
                    ),
                  ),
                ],
              ),
            ),
          if (searchQuery.isEmpty) const SizedBox(height: 16),

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
                        ],
                      ),
                    ),
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
                Text('Nenhum filme encontrado'),
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
          itemCount: movies.length,
          itemBuilder: (context, index) {
            return _buildSearchMovieCard(context, movies[index]);
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
      onTap: () =>
          context.go('/movie/${movie.id}?externalApiId=${movie.externalApiId}'),
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
    // Divide os filmes em diferentes carrosséis
    final trending = allMovies.take(10).toList();
    final popular = allMovies.skip(10).take(10).toList();
    final recommended = allMovies.skip(20).take(10).toList();

    return ListView(
      children: [
        _buildCarouselSection(context, 'Em Alta', trending),
        const SizedBox(height: 24),
        if (popular.isNotEmpty)
          _buildCarouselSection(context, 'Populares', popular),
        if (popular.isNotEmpty) const SizedBox(height: 24),
        if (recommended.isNotEmpty)
          _buildCarouselSection(context, 'Recomendados', recommended),
        if (recommended.isNotEmpty) const SizedBox(height: 24),
      ],
    );
  }

  Widget _buildCarouselSection(
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
            style: const TextStyle(fontSize: 20, fontWeight: FontWeight.bold),
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
              return _buildMovieCard(context, movies[index]);
            },
          ),
        ),
      ],
    );
  }

  Widget _buildMovieCard(BuildContext context, MovieModel movie) {
    return GestureDetector(
      onTap: () {
        context.go('/movie/${movie.id}?externalApiId=${movie.externalApiId}');
      },
      child: Container(
        width: 140,
        margin: const EdgeInsets.symmetric(horizontal: 4, vertical: 2),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          mainAxisSize: MainAxisSize.min,
          children: [
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
            Text(
              movie.title,
              maxLines: 1,
              overflow: TextOverflow.ellipsis,
              style: const TextStyle(fontWeight: FontWeight.w600, fontSize: 12),
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
      color: Colors.grey[800],
      child: const Center(
        child: Icon(Icons.movie, size: 50, color: Colors.grey),
      ),
    );
  }

  Widget _buildQuickAccessButton(
    BuildContext context,
    String title,
    IconData icon,
    String route,
  ) {
    return ElevatedButton(
      onPressed: () => context.go(route),
      style: ElevatedButton.styleFrom(
        padding: const EdgeInsets.symmetric(vertical: 12, horizontal: 16),
        backgroundColor: const Color(0xFFE50914).withOpacity(0.1),
        foregroundColor: const Color(0xFFE50914),
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(12),
          side: const BorderSide(color: Color(0xFFE50914), width: 1),
        ),
      ),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          Icon(icon, size: 20),
          const SizedBox(width: 8),
          Text(title, style: const TextStyle(fontWeight: FontWeight.w600)),
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
          const Padding(
            padding: EdgeInsets.symmetric(horizontal: 16.0, vertical: 8.0),
            child: Text(
              'Funcionalidades Futuras',
              style: TextStyle(
                fontSize: 12,
                fontWeight: FontWeight.w600,
                color: Colors.grey,
              ),
            ),
          ),
          ListTile(
            leading: const Icon(Icons.people, color: Colors.grey),
            title: const Text('Amigos', style: TextStyle(color: Colors.grey)),
            enabled: false,
          ),
          ListTile(
            leading: const Icon(Icons.favorite, color: Colors.grey),
            title: const Text(
              'Match de Filmes',
              style: TextStyle(color: Colors.grey),
            ),
            enabled: false,
          ),
          const Divider(),
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
