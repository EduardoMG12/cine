class UserModel {
  final String id;
  final String username;
  final String email;
  final String displayName;
  final String? bio;
  final String? profilePictureUrl;
  final bool isPrivate;
  final bool emailVerified;
  final String? theme;
  final DateTime createdAt;
  final DateTime updatedAt;

  UserModel({
    required this.id,
    required this.username,
    required this.email,
    required this.displayName,
    this.bio,
    this.profilePictureUrl,
    required this.isPrivate,
    required this.emailVerified,
    this.theme,
    required this.createdAt,
    required this.updatedAt,
  });

  factory UserModel.fromJson(Map<String, dynamic> json) {
    return UserModel(
      id: json['id'] as String,
      username: json['username'] as String,
      email: json['email'] as String,
      displayName: json['display_name'] as String,
      bio: json['bio'] as String?,
      profilePictureUrl: json['profile_picture_url'] as String?,
      isPrivate: json['is_private'] as bool? ?? false,
      emailVerified: json['email_verified'] as bool? ?? false,
      theme: json['theme'] as String?,
      createdAt: DateTime.parse(json['created_at'] as String),
      updatedAt: DateTime.parse(json['updated_at'] as String),
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'username': username,
      'email': email,
      'display_name': displayName,
      'bio': bio,
      'profile_picture_url': profilePictureUrl,
      'is_private': isPrivate,
      'email_verified': emailVerified,
      'theme': theme,
      'created_at': createdAt.toIso8601String(),
      'updated_at': updatedAt.toIso8601String(),
    };
  }
}
