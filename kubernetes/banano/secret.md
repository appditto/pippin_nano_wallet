// B64 encode wallet IDs
echo -n '["id", "id2"]' | base64

// Secret
apiVersion: v1
kind: Secret
metadata:
name: wallet-ids-banano
namespace: pippin
type: Opaque
data:
wallet-ids: b64encoded
