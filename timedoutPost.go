package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/pauldehodl/planet/x/blog/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	"strconv"
)

// GetTimedoutPostCount get the total number of timedoutPost
func (k Keeper) GetTimedoutPostCount(ctx sdk.Context) int64 {
	store :=  prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.TimedoutPostCountKey))
	byteKey := types.KeyPrefix(types.TimedoutPostCountKey)
	bz := store.Get(byteKey)

	// Count doesn't exist: no element
	if bz == nil {
		return 0
	}

	// Parse bytes
	count, err := strconv.ParseInt(string(bz), 10, 64)
	if err != nil {
		// Panic because the count should be always formattable to int64
		panic("cannot decode count")
	}

	return count
}

// SetTimedoutPostCount set the total number of timedoutPost
func (k Keeper) SetTimedoutPostCount(ctx sdk.Context, count int64)  {
	store :=  prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.TimedoutPostCountKey))
	byteKey := types.KeyPrefix(types.TimedoutPostCountKey)
	bz := []byte(strconv.FormatInt(count, 10))
	store.Set(byteKey, bz)
}

func (k Keeper) CreateTimedoutPost(ctx sdk.Context, msg types.MsgCreateTimedoutPost) {
	// Create the timedoutPost
    count := k.GetTimedoutPostCount(ctx)
    var timedoutPost = types.TimedoutPost{
        Creator: msg.Creator,
        Id:      strconv.FormatInt(count, 10),
        Title: msg.Title,
        Chain: msg.Chain,
    }

    store :=  prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.TimedoutPostKey))
    key := types.KeyPrefix(types.TimedoutPostKey + timedoutPost.Id)
    value := k.cdc.MustMarshalBinaryBare(&timedoutPost)
    store.Set(key, value)

    // Update timedoutPost count
    k.SetTimedoutPostCount(ctx, count+1)
}

func (k Keeper) UpdateTimedoutPost(ctx sdk.Context, timedoutPost types.TimedoutPost) {
	store :=  prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.TimedoutPostKey))
	b := k.cdc.MustMarshalBinaryBare(&timedoutPost)
	store.Set(types.KeyPrefix(types.TimedoutPostKey + timedoutPost.Id), b)
}

func (k Keeper) GetTimedoutPost(ctx sdk.Context, key string) types.TimedoutPost {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.TimedoutPostKey))
	var timedoutPost types.TimedoutPost
	k.cdc.MustUnmarshalBinaryBare(store.Get(types.KeyPrefix(types.TimedoutPostKey + key)), &timedoutPost)
	return timedoutPost
}

func (k Keeper) HasTimedoutPost(ctx sdk.Context, id string) bool {
	store :=  prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.TimedoutPostKey))
	return store.Has(types.KeyPrefix(types.TimedoutPostKey + id))
}

func (k Keeper) GetTimedoutPostOwner(ctx sdk.Context, key string) string {
    return k.GetTimedoutPost(ctx, key).Creator
}

// DeleteTimedoutPost deletes a timedoutPost
func (k Keeper) DeleteTimedoutPost(ctx sdk.Context, key string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.TimedoutPostKey))
	store.Delete(types.KeyPrefix(types.TimedoutPostKey + key))
}

func (k Keeper) GetAllTimedoutPost(ctx sdk.Context) (msgs []types.TimedoutPost) {
    store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.TimedoutPostKey))
	iterator := sdk.KVStorePrefixIterator(store, types.KeyPrefix(types.TimedoutPostKey))

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var msg types.TimedoutPost
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &msg)
        msgs = append(msgs, msg)
	}

    return
}
