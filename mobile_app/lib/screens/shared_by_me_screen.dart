import 'package:flutter/material.dart';
import '../services/api_service.dart';

class SharedByMeScreen extends StatefulWidget {
  final String token;

  const SharedByMeScreen({Key? key, required this.token}) : super(key: key);

  @override
  _SharedByMeScreenState createState() => _SharedByMeScreenState();
}

class _SharedByMeScreenState extends State<SharedByMeScreen> {
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
      final documents = await _apiService.getSharedByMe(widget.token);
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
        title: const Text('Shared By Me'),
      ),
      body: _isLoading
          ? const Center(child: CircularProgressIndicator())
          : _sharedDocuments.isEmpty
              ? const Center(child: Text('No shared documents'))
              : ListView.builder(
                  itemCount: _sharedDocuments.length,
                  itemBuilder: (context, index) {
                    final share = _sharedDocuments[index];
                    return ListTile(
                      title: Text(share['document']['title'] ?? 'Untitled'),
                      subtitle: Text(
                        'Shared with: ${share['user']['username']} (${share['user']['email']})',
                      ),
                      trailing: IconButton(
                        icon: const Icon(Icons.delete, color: Colors.red),
                        onPressed: () async {
                          try {
                            await _apiService.removeShare(
                              int.parse(share['id'].toString()),
                              widget.token,
                            );
                            setState(() {
                              _sharedDocuments.removeAt(index);
                            });
                            ScaffoldMessenger.of(context).showSnackBar(
                              const SnackBar(content: Text('Share removed')),
                            );
                          } catch (e) {
                            ScaffoldMessenger.of(context).showSnackBar(
                              SnackBar(
                                  content: Text('Failed to remove share: $e')),
                            );
                          }
                        },
                      ),
                    );
                  },
                ),
    );
  }
}
