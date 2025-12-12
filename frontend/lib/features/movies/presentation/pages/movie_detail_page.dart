import 'package:flutter/material.dart';

class MovieDetailPage extends StatelessWidget {
  final String movieId;
  
  const MovieDetailPage({super.key, required this.movieId});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Movie Details')),
      body: Center(child: Text('Movie ID: $movieId')),
    );
  }
}

class WatchLaterPage extends StatelessWidget {
  const WatchLaterPage({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Watch Later')),
      body: const Center(child: Text('Watch Later List')),
    );
  }
}

class WatchedMoviesPage extends StatelessWidget {
  const WatchedMoviesPage({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Watched Movies')),
      body: const Center(child: Text('Watched Movies List')),
    );
  }
}
