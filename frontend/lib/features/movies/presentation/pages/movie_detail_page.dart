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
