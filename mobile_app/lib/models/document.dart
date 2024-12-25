class Document {
  final int id;
  final String title;
  final String content;
  final bool isPublic;
  final int userId;
  final List<String> tags;

  Document({
    required this.id,
    required this.title,
    required this.content,
    required this.isPublic,
    required this.userId,
    required this.tags,
  });

  factory Document.fromJson(Map<String, dynamic> json) {
    return Document(
      id: json['id'],
      title: json['title'],
      content: json['content'],
      isPublic: json['is_public'],
      userId: json['user_id'],
      tags: List<String>.from(json['tags'].map((x) => x)),
    );
  }
}
