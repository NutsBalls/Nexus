import 'package:flutter/material.dart';
import '../services/api_service.dart';
import '../screens/document_detail_screen.dart';

class SharedWithMeScreen extends StatefulWidget {
  final String token;

  const SharedWithMeScreen({Key? key, required this.token}) : super(key: key);

  @override
  _SharedWithMeScreenState createState() => _SharedWithMeScreenState();
}

class _SharedWithMeScreenState extends State<SharedWithMeScreen> {
  final ApiService _apiService = ApiService();
  List<Map<String, dynamic>> _sharedDocuments = [];
  bool _isLoading = true;

  @override
  void initState() {
    super.initState();
    _loadSharedDocuments();
  }

  Future<void> _loadSharedDocuments() async {
    setState(() => _isLoading = true);
    try {
      final documents = await _apiService.getSharedWithMe(widget.token);
      setState(() {
        _sharedDocuments = documents;
        _isLoading = false;
      });
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('Failed to load shared documents: $e')),
        );
      }
      setState(() => _isLoading = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Shared With Me'),
      ),
      body: _isLoading
          ? const Center(child: CircularProgressIndicator())
          : _sharedDocuments.isEmpty
              ? const Center(child: Text('No shared documents'))
              : ListView.builder(
                  itemCount: _sharedDocuments.length,
                  itemBuilder: (context, index) {
                    final document = _sharedDocuments[index];
                    return ListTile(
                      title: Text(document['title'] ?? 'Untitled'),
                      subtitle: Text('Shared by: ${document['user_id']}'),
                      onTap: () {
                        Navigator.push(
                          context,
                          MaterialPageRoute(
                            builder: (context) => DocumentDetailScreen(
                              token: widget.token,
                              id: document['id'],
                            ),
                          ),
                        );
                      },
                    );
                  },
                ),
    );
  }
}
