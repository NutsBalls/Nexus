import 'package:flutter/material.dart';
import '../services/api_service.dart';

class ShareScreen extends StatefulWidget {
  final String token;
  final int documentId;
  final String documentTitle;

  const ShareScreen({
    Key? key,
    required this.token,
    required this.documentId,
    required this.documentTitle,
  }) : super(key: key);

  @override
  _ShareScreenState createState() => _ShareScreenState();
}

class _ShareScreenState extends State<ShareScreen> {
  final _emailController = TextEditingController();
  final ApiService _apiService = ApiService();
  List<Map<String, dynamic>> shares = [];
  bool _isLoading = true;

  @override
  void initState() {
    super.initState();
    _loadShares();
  }

  Future<void> _loadShares() async {
    try {
      final response = await _apiService.getDocumentShares(
        widget.documentId,
        widget.token,
      );
      setState(() {
        shares = List<Map<String, dynamic>>.from(response);
        _isLoading = false;
      });
    } catch (e) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('Failed to load shares: $e')),
      );
      setState(() => _isLoading = false);
    }
  }

  Future<void> _shareDocument() async {
    final permission = await showDialog<String>(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Select Permission'),
        content: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            ListTile(
              title: const Text('Read Only'),
              onTap: () => Navigator.pop(context, 'read'),
            ),
            ListTile(
              title: const Text('Can Edit'),
              onTap: () => Navigator.pop(context, 'write'),
            ),
            ListTile(
              title: const Text('Admin'),
              onTap: () => Navigator.pop(context, 'admin'),
            ),
          ],
        ),
      ),
    );

    if (permission != null && _emailController.text.isNotEmpty) {
      print('Email before sending: ${_emailController.text}');
      try {
        await _apiService.shareDocument(
          widget.documentId,
          _emailController.text,
          permission,
          widget.token,
        );
        print('Email after sending: ${_emailController.text}');
        _emailController.clear();
        _loadShares();
        if (mounted) {
          ScaffoldMessenger.of(context).showSnackBar(
            const SnackBar(content: Text('Document shared successfully')),
          );
        }
      } catch (e) {
        if (mounted) {
          ScaffoldMessenger.of(context).showSnackBar(
            SnackBar(content: Text('Failed to share document: $e')),
          );
        }
      }
    } else {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(content: Text('Please provide an email.')),
        );
      }
    }
  }

  Future<void> _removeShare(int shareId) async {
    final confirmed = await showDialog<bool>(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Remove Access'),
        content:
            const Text('Are you sure you want to remove access for this user?'),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context, false),
            child: const Text('Cancel'),
          ),
          TextButton(
            onPressed: () => Navigator.pop(context, true),
            style: TextButton.styleFrom(foregroundColor: Colors.red),
            child: const Text('Remove'),
          ),
        ],
      ),
    );

    if (confirmed == true) {
      try {
        await _apiService.removeShare(shareId, widget.token);
        _loadShares();
        if (mounted) {
          ScaffoldMessenger.of(context).showSnackBar(
            const SnackBar(content: Text('Access removed successfully')),
          );
        }
      } catch (e) {
        if (mounted) {
          ScaffoldMessenger.of(context).showSnackBar(
            SnackBar(content: Text('Failed to remove access: $e')),
          );
        }
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text('Share ${widget.documentTitle}'),
      ),
      body: Column(
        children: [
          Padding(
            padding: const EdgeInsets.all(16.0),
            child: Row(
              children: [
                Expanded(
                  child: TextField(
                    controller: _emailController,
                    onChanged: (text) {
                      print('Text field changed: $text');
                    },
                    decoration: const InputDecoration(
                      labelText: 'User Email',
                      hintText: 'Enter email to share with',
                      border: OutlineInputBorder(),
                    ),
                  ),
                ),
                const SizedBox(width: 16),
                ElevatedButton(
                  onPressed: _shareDocument,
                  child: const Text('Share'),
                ),
              ],
            ),
          ),
          Expanded(
            child: _isLoading
                ? const Center(child: CircularProgressIndicator())
                : shares.isEmpty
                    ? const Center(child: Text('No shares yet'))
                    : ListView.builder(
                        itemCount: shares.length,
                        itemBuilder: (context, index) {
                          final share = shares[index];
                          final user = share['user'];
                          return ListTile(
                            leading: CircleAvatar(
                              child: Text(user['username'][0].toUpperCase()),
                            ),
                            title: Text(user['email']),
                            subtitle: Text(
                              'Permission: ${share['permission']}',
                              style: TextStyle(
                                color: _getPermissionColor(share['permission']),
                              ),
                            ),
                            trailing: IconButton(
                              icon: const Icon(Icons.delete),
                              onPressed: () => _removeShare(share['id']),
                            ),
                          );
                        },
                      ),
          ),
        ],
      ),
    );
  }

  Color _getPermissionColor(String permission) {
    switch (permission) {
      case 'read':
        return Colors.blue;
      case 'write':
        return Colors.green;
      case 'admin':
        return Colors.orange;
      default:
        return Colors.grey;
    }
  }

  @override
  void dispose() {
    _emailController.dispose();
    super.dispose();
  }
}
