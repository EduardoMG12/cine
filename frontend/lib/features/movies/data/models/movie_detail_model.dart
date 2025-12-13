class MovieRating {
  final String source;
  final String value;

  MovieRating({
    required this.source,
    required this.value,
  });

  factory MovieRating.fromJson(Map<String, dynamic> json) {
    return MovieRating(
      source: json['source'] ?? '',
      value: json['value'] ?? '',
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'source': source,
      'value': value,
    };
  }
}

class MovieDetailModel {
  final String title;
  final String year;
  final String released;
  final String runtime;
  final String plot;
  final String type;
  final String poster;
  final String rated;
  final String genre;
  final String language;
  final String country;
  final String director;
  final String writer;
  final String actors;
  final String imdbId;
  final String imdbRating;
  final String imdbVotes;
  final String metascore;
  final List<MovieRating> ratings;
  final String awards;
  final String? boxOffice;
  final String? production;
  final String? website;
  final String provider;
  final String providerId;

  MovieDetailModel({
    required this.title,
    required this.year,
    required this.released,
    required this.runtime,
    required this.plot,
    required this.type,
    required this.poster,
    required this.rated,
    required this.genre,
    required this.language,
    required this.country,
    required this.director,
    required this.writer,
    required this.actors,
    required this.imdbId,
    required this.imdbRating,
    required this.imdbVotes,
    required this.metascore,
    required this.ratings,
    required this.awards,
    this.boxOffice,
    this.production,
    this.website,
    required this.provider,
    required this.providerId,
  });

  factory MovieDetailModel.fromJson(Map<String, dynamic> json) {
    return MovieDetailModel(
      title: json['title'] ?? '',
      year: json['year'] ?? '',
      released: json['released'] ?? '',
      runtime: json['runtime'] ?? '',
      plot: json['plot'] ?? '',
      type: json['type'] ?? '',
      poster: json['poster'] ?? '',
      rated: json['rated'] ?? '',
      genre: json['genre'] ?? '',
      language: json['language'] ?? '',
      country: json['country'] ?? '',
      director: json['director'] ?? '',
      writer: json['writer'] ?? '',
      actors: json['actors'] ?? '',
      imdbId: json['imdb_id'] ?? '',
      imdbRating: json['imdb_rating'] ?? '',
      imdbVotes: json['imdb_votes'] ?? '',
      metascore: json['metascore'] ?? '',
      ratings: (json['ratings'] as List<dynamic>?)
              ?.map((r) => MovieRating.fromJson(r))
              .toList() ??
          [],
      awards: json['awards'] ?? '',
      boxOffice: json['box_office'],
      production: json['production'],
      website: json['website'],
      provider: json['provider'] ?? '',
      providerId: json['provider_id'] ?? '',
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'title': title,
      'year': year,
      'released': released,
      'runtime': runtime,
      'plot': plot,
      'type': type,
      'poster': poster,
      'rated': rated,
      'genre': genre,
      'language': language,
      'country': country,
      'director': director,
      'writer': writer,
      'actors': actors,
      'imdb_id': imdbId,
      'imdb_rating': imdbRating,
      'imdb_votes': imdbVotes,
      'metascore': metascore,
      'ratings': ratings.map((r) => r.toJson()).toList(),
      'awards': awards,
      'box_office': boxOffice,
      'production': production,
      'website': website,
      'provider': provider,
      'provider_id': providerId,
    };
  }

  List<String> get genresList {
    return genre.split(',').map((g) => g.trim()).toList();
  }
}
