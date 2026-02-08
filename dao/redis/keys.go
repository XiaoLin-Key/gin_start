package redis

const (
	KeyPrefix          = "bluebell:"
	KeyPostTimeZSet    = "post:time"
	KeyPostScoreZSet   = "post:score"
	KeyPostVotedZSetPF = "post:voted:"

	KeyCommunitySetPF = "community:"
)

func getRedisKey(key string) string {
	return KeyPrefix + key
}
