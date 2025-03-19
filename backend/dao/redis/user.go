package redis

import "go.uber.org/zap"

func UserTokenNum(uid string) (int, error) {
	key := KeyUserTokenSetPrefix + uid + KeyUserTokenSetSuffix
	n, err := cli.SCard(key).Result()
	if err != nil {
		zap.L().Error("cli.SCard error", zap.Error(err))
		return -2, err
	}

	return int(n), err
}

func UserTokenList(uid string) ([]string, error) {
	key := KeyUserTokenSetPrefix + uid + KeyUserTokenSetSuffix
	tokens, err := cli.SMembers(key).Result()
	if err != nil {
		zap.L().Error("cli.SMembers error", zap.Error(err))
		return nil, err
	}
	return tokens, err
}

func SubUserToken(uid string) error {
	key := KeyUserTokenSetPrefix + uid + KeyUserTokenSetSuffix

	_, err := cli.SPop(key).Result()
	if err != nil {
		zap.L().Error("cli.SPop error", zap.Error(err))
		return err
	}
	return nil

}

func AddUserToken(uid, token string) error {
	key := KeyUserTokenSetPrefix + uid + KeyUserTokenSetSuffix

	if err := cli.SAdd(key, token).Err(); err != nil {
		zap.L().Error("cli.SAdd error", zap.Error(err))
		return err
	}

	return nil
}

func Logout(uid string, token string) error {
	key := KeyUserTokenSetPrefix + uid + KeyUserTokenSetSuffix

	if err := cli.SRem(key, token).Err(); err != nil {
		zap.L().Error("cli.SRem error", zap.Error(err))
		return err
	}
	return nil
}
