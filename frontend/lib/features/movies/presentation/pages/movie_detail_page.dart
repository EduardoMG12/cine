import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import '../providers/movie_detail_provider.dart';
import '../../data/models/movie_detail_model.dart';
import '../../../auth/presentation/providers/auth_provider.dart';
import '../../../../core/services/user_movie_service.dart';

class MovieDetailPage extends ConsumerStatefulWidget {
  final String movieId; // UUID do filme no banco de dados
  final String? externalApiId; // ID da API externa (OMDB/TMDB)

  const MovieDetailPage({super.key, required this.movieId, this.externalApiId});

  @override
  ConsumerState<MovieDetailPage> createState() => _MovieDetailPageState();
}

class _MovieDetailPageState extends ConsumerState<MovieDetailPage> {
  bool _isTogglingWatched = false;
  bool _isTogglingFavorite = false;

  Future<void> _toggleWatched(BuildContext context) async {
    setState(() {
      _isTogglingWatched = true;
    });

    try {
      final result = await UserMovieService.toggleWatchedMovie(widget.movieId);

      if (context.mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text(result['message'] ?? 'Status atualizado'),
            backgroundColor: result['added'] == true
                ? Colors.green
                : Colors.orange,
          ),
        );
      }
    } catch (e) {
      if (context.mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text('Erro ao atualizar: $e'),
            backgroundColor: Colors.red,
          ),
        );
      }
    } finally {
      if (mounted) {
        setState(() {
          _isTogglingWatched = false;
        });
      }
    }
  }

  Future<void> _toggleFavorite(BuildContext context) async {
    setState(() {
      _isTogglingFavorite = true;
    });

    try {
      final result = await UserMovieService.toggleFavoriteMovie(widget.movieId);

      if (context.mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text(result['message'] ?? 'Status atualizado'),
            backgroundColor: result['added'] == true
                ? Colors.green
                : Colors.orange,
          ),
        );
      }
    } catch (e) {
      if (context.mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text('Erro ao atualizar: $e'),
            backgroundColor: Colors.red,
          ),
        );
      }
    } finally {
      if (mounted) {
        setState(() {
          _isTogglingFavorite = false;
        });
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    // Usa externalApiId para buscar detalhes, ou fallback para movieId
    final apiId = widget.externalApiId ?? widget.movieId;
    final movieDetailAsync = ref.watch(movieDetailProvider(apiId));
    final authState = ref.watch(authStateProvider);

    return Scaffold(
      body: movieDetailAsync.when(
        data: (movieDetail) => _buildMovieDetail(
          context,
          ref,
          movieDetail,
          authState.isAuthenticated,
        ),
        loading: () => const Center(child: CircularProgressIndicator()),
        error: (error, stack) => Scaffold(
          appBar: AppBar(title: const Text('Erro')),
          body: Center(
            child: Column(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                const Icon(Icons.error_outline, size: 48, color: Colors.red),
                const SizedBox(height: 16),
                Text('Erro ao carregar detalhes: $error'),
                const SizedBox(height: 16),
                ElevatedButton(
                  onPressed: () =>
                      ref.refresh(movieDetailProvider(widget.movieId)),
                  child: const Text('Tentar Novamente'),
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }

  Widget _buildMovieDetail(
    BuildContext context,
    WidgetRef ref,
    MovieDetailModel movie,
    bool isLoggedIn,
  ) {
    return CustomScrollView(
      slivers: [
        // App Bar com poster de fundo
        SliverAppBar(
          expandedHeight: 400,
          pinned: true,
          leading: IconButton(
            icon: const Icon(Icons.arrow_back),
            onPressed: () {
              // Se conseguir voltar usando o Navigator, usa pop
              if (Navigator.of(context).canPop()) {
                Navigator.of(context).pop();
              } else {
                // Caso contrário, redireciona para a home apropriada
                if (isLoggedIn) {
                  context.go('/home');
                } else {
                  context.go('/');
                }
              }
            },
          ),
          flexibleSpace: FlexibleSpaceBar(
            title: Text(
              movie.title,
              style: const TextStyle(
                shadows: [Shadow(color: Colors.black, blurRadius: 10)],
              ),
            ),
            background: Stack(
              fit: StackFit.expand,
              children: [
                // Poster
                movie.poster.isNotEmpty
                    ? Image.network(
                        movie.poster,
                        fit: BoxFit.cover,
                        errorBuilder: (context, error, stackTrace) {
                          return _buildPlaceholder();
                        },
                      )
                    : _buildPlaceholder(),
                // Gradiente overlay
                Container(
                  decoration: BoxDecoration(
                    gradient: LinearGradient(
                      begin: Alignment.topCenter,
                      end: Alignment.bottomCenter,
                      colors: [
                        Colors.transparent,
                        Colors.black.withOpacity(0.7),
                      ],
                    ),
                  ),
                ),
              ],
            ),
          ),
        ),

        // Conteúdo
        SliverToBoxAdapter(
          child: Padding(
            padding: const EdgeInsets.all(16),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                // Informações básicas
                Row(
                  children: [
                    _buildInfoChip(movie.year),
                    const SizedBox(width: 8),
                    _buildInfoChip(movie.runtime),
                    const SizedBox(width: 8),
                    _buildInfoChip(movie.rated),
                  ],
                ),

                const SizedBox(height: 16),

                // Gêneros
                Wrap(
                  spacing: 8,
                  runSpacing: 8,
                  children: movie.genresList.map((genre) {
                    return Chip(
                      label: Text(genre),
                      backgroundColor: const Color(0xFFE50914).withOpacity(0.2),
                    );
                  }).toList(),
                ),

                // Action Buttons (Watched and Favorite) - Only for authenticated users
                if (isLoggedIn) ...[
                  const SizedBox(height: 24),
                  Row(
                    children: [
                      Expanded(
                        child: ElevatedButton.icon(
                          onPressed: _isTogglingWatched
                              ? null
                              : () => _toggleWatched(context),
                          icon: _isTogglingWatched
                              ? const SizedBox(
                                  width: 20,
                                  height: 20,
                                  child: CircularProgressIndicator(
                                    strokeWidth: 2,
                                  ),
                                )
                              : const Icon(Icons.check_circle_outline),
                          label: const Text('Assistido'),
                          style: ElevatedButton.styleFrom(
                            backgroundColor: Colors.green.withOpacity(0.2),
                            foregroundColor: Colors.green,
                            padding: const EdgeInsets.symmetric(vertical: 12),
                          ),
                        ),
                      ),
                      const SizedBox(width: 12),
                      Expanded(
                        child: ElevatedButton.icon(
                          onPressed: _isTogglingFavorite
                              ? null
                              : () => _toggleFavorite(context),
                          icon: _isTogglingFavorite
                              ? const SizedBox(
                                  width: 20,
                                  height: 20,
                                  child: CircularProgressIndicator(
                                    strokeWidth: 2,
                                  ),
                                )
                              : const Icon(Icons.favorite_border),
                          label: const Text('Favorito'),
                          style: ElevatedButton.styleFrom(
                            backgroundColor: const Color(
                              0xFFE50914,
                            ).withOpacity(0.2),
                            foregroundColor: const Color(0xFFE50914),
                            padding: const EdgeInsets.symmetric(vertical: 12),
                          ),
                        ),
                      ),
                    ],
                  ),
                ],

                const SizedBox(height: 24),

                // Plot
                Text(
                  'Sinopse',
                  style: Theme.of(
                    context,
                  ).textTheme.titleLarge?.copyWith(fontWeight: FontWeight.bold),
                ),
                const SizedBox(height: 8),
                Text(movie.plot, style: Theme.of(context).textTheme.bodyLarge),

                const SizedBox(height: 24),

                // Ratings
                Text(
                  'Avaliações',
                  style: Theme.of(
                    context,
                  ).textTheme.titleLarge?.copyWith(fontWeight: FontWeight.bold),
                ),
                const SizedBox(height: 12),
                ...movie.ratings.map((rating) {
                  return Padding(
                    padding: const EdgeInsets.only(bottom: 8),
                    child: Row(
                      mainAxisAlignment: MainAxisAlignment.spaceBetween,
                      children: [
                        Text(
                          rating.source,
                          style: TextStyle(color: Colors.grey[400]),
                        ),
                        Text(
                          rating.value,
                          style: const TextStyle(
                            fontWeight: FontWeight.bold,
                            fontSize: 16,
                          ),
                        ),
                      ],
                    ),
                  );
                }),

                const SizedBox(height: 24),

                // IMDb Rating destacado
                Container(
                  padding: const EdgeInsets.all(16),
                  decoration: BoxDecoration(
                    color: const Color(0xFFF5C518).withOpacity(0.2),
                    borderRadius: BorderRadius.circular(12),
                    border: Border.all(
                      color: const Color(0xFFF5C518),
                      width: 2,
                    ),
                  ),
                  child: Row(
                    mainAxisAlignment: MainAxisAlignment.spaceAround,
                    children: [
                      Column(
                        children: [
                          const Text(
                            'IMDb Rating',
                            style: TextStyle(
                              color: Color(0xFFF5C518),
                              fontWeight: FontWeight.bold,
                            ),
                          ),
                          const SizedBox(height: 4),
                          Text(
                            movie.imdbRating,
                            style: const TextStyle(
                              fontSize: 24,
                              fontWeight: FontWeight.bold,
                            ),
                          ),
                          Text(
                            '${movie.imdbVotes} votos',
                            style: TextStyle(
                              fontSize: 12,
                              color: Colors.grey[400],
                            ),
                          ),
                        ],
                      ),
                      if (movie.metascore.isNotEmpty &&
                          movie.metascore != 'N/A')
                        Column(
                          children: [
                            const Text(
                              'Metascore',
                              style: TextStyle(
                                color: Color(0xFF66CC33),
                                fontWeight: FontWeight.bold,
                              ),
                            ),
                            const SizedBox(height: 4),
                            Text(
                              movie.metascore,
                              style: const TextStyle(
                                fontSize: 24,
                                fontWeight: FontWeight.bold,
                              ),
                            ),
                          ],
                        ),
                    ],
                  ),
                ),

                const SizedBox(height: 24),

                // Elenco e Equipe
                _buildInfoSection(context, 'Diretor', movie.director),
                _buildInfoSection(context, 'Roteirista', movie.writer),
                _buildInfoSection(context, 'Elenco', movie.actors),

                const SizedBox(height: 24),

                // Informações adicionais
                _buildInfoSection(context, 'Idioma', movie.language),
                _buildInfoSection(context, 'País', movie.country),
                _buildInfoSection(context, 'Lançamento', movie.released),

                if (movie.awards.isNotEmpty && movie.awards != 'N/A') ...[
                  const SizedBox(height: 16),
                  _buildInfoSection(context, 'Prêmios', movie.awards),
                ],

                const SizedBox(height: 32),
              ],
            ),
          ),
        ),
      ],
    );
  }

  Widget _buildInfoChip(String text) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 6),
      decoration: BoxDecoration(
        color: Colors.grey[800],
        borderRadius: BorderRadius.circular(20),
      ),
      child: Text(text, style: const TextStyle(fontSize: 12)),
    );
  }

  Widget _buildInfoSection(BuildContext context, String label, String value) {
    if (value.isEmpty || value == 'N/A') return const SizedBox.shrink();

    return Padding(
      padding: const EdgeInsets.only(bottom: 12),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            label,
            style: TextStyle(
              color: Colors.grey[400],
              fontSize: 14,
              fontWeight: FontWeight.w600,
            ),
          ),
          const SizedBox(height: 4),
          Text(value, style: const TextStyle(fontSize: 16)),
        ],
      ),
    );
  }

  Widget _buildPlaceholder() {
    return Container(
      color: Colors.grey[800],
      child: const Center(
        child: Icon(Icons.movie, size: 100, color: Colors.grey),
      ),
    );
  }
}
