import 'dart:math' as math;
import 'package:flutter/material.dart';
import '../services/api_service.dart';

class DocumentsScreen extends StatefulWidget {
  final String token;
  final int folderId;
  final String folderName;

  const DocumentsScreen({
    Key? key,
    required this.token,
    required this.folderId,
    required this.folderName,
  }) : super(key: key);

  @override
  _DocumentsScreenState createState() => _DocumentsScreenState();
}

class _DocumentsScreenState extends State<DocumentsScreen> {
  final ApiService _apiService = ApiService();
  List<Map<String, dynamic>> documents = [];
  bool _isLoading = true;

  @override
  void initState() {
    super.initState();
    _fetchDocuments();
  }

  Future<void> _fetchDocuments() async {
    setState(() => _isLoading = true);
    try {
      final response = await _apiService.get(
        'folders/${widget.folderId}/documents',
        widget.token,
      );
      setState(() {
        documents = List<Map<String, dynamic>>.from(response);
        _isLoading = false;
      });
    } catch (e) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('Failed to load documents: $e')),
      );
      setState(() => _isLoading = false);
    }
  }

  Future<void> _deleteDocument(int documentId) async {
    final confirmed = await showDialog<bool>(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Delete Document'),
        content: const Text(
          'Are you sure you want to delete this document? '
          'This action cannot be undone.',
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context, false),
            child: const Text('Cancel'),
          ),
          TextButton(
            onPressed: () => Navigator.pop(context, true),
            style: TextButton.styleFrom(
              foregroundColor: Colors.red,
            ),
            child: const Text('Delete'),
          ),
        ],
      ),
    );

    if (confirmed == true) {
      try {
        await _apiService.deleteDocument(documentId, widget.token);
        _fetchDocuments();
        if (mounted) {
          ScaffoldMessenger.of(context).showSnackBar(
            const SnackBar(content: Text('Document deleted successfully')),
          );
        }
      } catch (e) {
        if (mounted) {
          ScaffoldMessenger.of(context).showSnackBar(
            SnackBar(content: Text('Failed to delete document: $e')),
          );
        }
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text(widget.folderName),
        leading: IconButton(
          icon: const Icon(Icons.arrow_back),
          onPressed: () => Navigator.pop(context),
        ),
      ),
      body: _isLoading
          ? const Center(child: CircularProgressIndicator())
          : documents.isEmpty
              ? const Center(child: Text('No documents in this folder'))
              : ListView.builder(
                  itemCount: documents.length,
                  itemBuilder: (context, index) {
                    final document = documents[index];
                    return Dismissible(
                      key: Key(document['id'].toString()),
                      direction: DismissDirection.endToStart,
                      background: Container(
                        color: Colors.red,
                        alignment: Alignment.centerRight,
                        padding: const EdgeInsets.only(right: 16),
                        child: const Icon(
                          Icons.delete,
                          color: Colors.white,
                        ),
                      ),
                      confirmDismiss: (direction) async {
                        return await showDialog<bool>(
                          context: context,
                          builder: (context) => AlertDialog(
                            title: const Text('Delete Document'),
                            content: const Text(
                              'Are you sure you want to delete this document? '
                              'This action cannot be undone.',
                            ),
                            actions: [
                              TextButton(
                                onPressed: () => Navigator.pop(context, false),
                                child: const Text('Cancel'),
                              ),
                              TextButton(
                                onPressed: () => Navigator.pop(context, true),
                                style: TextButton.styleFrom(
                                  foregroundColor: Colors.red,
                                ),
                                child: const Text('Delete'),
                              ),
                            ],
                          ),
                        );
                      },
                      onDismissed: (direction) {
                        _deleteDocument(document['id']);
                      },
                      child: ListTile(
                        leading: const Icon(Icons.description),
                        title: Text(document['title'] ?? 'Untitled'),
                        subtitle: Text(
                          document['content']?.toString().substring(
                                  0,
                                  math.min(
                                      50,
                                      document['content']?.toString().length ??
                                          0)) ??
                              'No content',
                        ),
                        trailing: PopupMenuButton(
                          itemBuilder: (context) => [
                            const PopupMenuItem(
                              value: 'delete',
                              child: Text('Delete'),
                            ),
                            const PopupMenuItem(
                              value: 'edit',
                              child: Text('Edit'),
                            ),
                          ],
                          onSelected: (value) {
                            switch (value) {
                              case 'delete':
                                _deleteDocument(document['id']);
                                break;
                              case 'edit':
                                break;
                            }
                          },
                        ),
                        onTap: () {
                          Navigator.pushNamed(
                            context,
                            '/document-detail',
                            arguments: {
                              'id': document['id'],
                              'token': widget.token,
                            },
                          ).then((_) => _fetchDocuments());
                        },
                      ),
                    );
                  },
                ),
      floatingActionButton: FloatingActionButton(
        onPressed: () {
          Navigator.pushNamed(
            context,
            '/create-document',
            arguments: {
              'token': widget.token,
              'folderId': widget.folderId,
              'folderName': widget.folderName,
            },
          ).then((_) => _fetchDocuments());
        },
        child: const Icon(Icons.add),
      ),
    );
  }
}
