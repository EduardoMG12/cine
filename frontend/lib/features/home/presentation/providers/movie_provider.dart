import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../../../core/services/movie_service.dart';
import '../../../movies/data/models/movie_model.dart';

/// Provider para buscar filmes trending
final trendingMoviesProvider = FutureProvider<List<MovieModel>>((ref) async {
  return await MovieService.getTrendingMovies();
});

/// Notifier para o termo de pesquisa
class SearchQueryNotifier extends Notifier<String> {
  @override
  String build() => '';

  void setQuery(String query) {
    state = query;
  }
}

final searchQueryProvider = NotifierProvider<SearchQueryNotifier, String>(() {
  return SearchQueryNotifier();
});

/// Notifier para a página atual de pesquisa
class SearchPageNotifier extends Notifier<int> {
  @override
  int build() => 1;

  void nextPage() {
    state = state + 1;
  }

  void previousPage() {
    if (state > 1) {
      state = state - 1;
    }
  }

  void reset() {
    state = 1;
  }
}

final searchPageProvider = NotifierProvider<SearchPageNotifier, int>(() {
  return SearchPageNotifier();
});

/// Provider para buscar filmes por pesquisa
final searchMoviesProvider =
    FutureProvider.family<List<MovieModel>, ({String query, int page})>((
      ref,
      params,
    ) async {
      if (params.query.isEmpty) {
        return [];
      }

      return await MovieService.searchMovies(params.query, page: params.page);
    });
