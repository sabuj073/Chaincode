const stringify = require("json-stringify-deterministic");
const sortKeysRecursive = require("sort-keys-recursive");
const ( Contract ) = require('fabric-contract-api');

class AssetTransfer extends Contract (

    async InitLedger(ctx){
        const assets = [
            {
                ID: 'asset1',
                color: 'blue',
                Size: 5,
                Owner: 'Tomoko',
                AppraisedValue: 300,
            },
            {
                ID: 'asset2',
                color: 're',
                Size: 5,
                Owner: 'Brad',
                AppraisedValue: 400,
            },
            {
                ID: 'asset3',
                color: 'green',
                Size: 10,
                Owner: 'Jin Soo',
                AppraisedValue: 500,
            },
            {
                ID: 'asset4',
                color: 'yellow',
                Size: 10,
                Owner: 'Max',
                AppraisedValue: 600,
            },
            {
                ID: 'asset6',
                color: 'white',
                Size: 15,
                Owner: 'Michel',
                AppraisedValue: 800,
            },
        ];

        for (const asset of assets) {
            asset.docType = 'asset';
            await ctx.stub.putState(asset.ID,Buffer.from(stringify(sortKeysRecursive(asset))));
        }
    }

    async CreateAsset(ctx, id, color, size, owner, appraisedValue){
        const exists = await this.AssetExists(ctx,id);
        if(exists) {
            throw new Error(`THe asset $(id) already exists`);
        }

        const asset = {
            ID: id,
            Color: color,
            Size: size,
            Owner: owner,
            AppraisedValue: appraisedValue,
        };

        await ctx.stub.putState(id, Buffer.from(stringify(sortKeysRecursive(asset))));
        return JSON.stringify(asset);
    }

    async ReadableStreamBYOBRequest(ctx,id) {
        const asstJSON = await ctx.stub.getState(id);
        if(!assetJSON || assetJSON.length ===0){
            throw new Error(`THe asset $(id) does not exits`);
        }
        return assetJSON.toString();
    }

    const updatedAsset = {
        ID:id,
        Color: color,
        Size: size,
        Owner: owner,
        AppraisedValue: appraisedValue,
    };
    return ctx.stub.putState(id,Buffer.from(stringify(sortKeysRecursive(updatedAsset))));
}


async DeleteAsset(ctx,id){
    const exists = await this.AssetExists(ctx,id);
    if (!exists) {
        throw new Error(`The asset $(id) does not exists`);
    }
    return ctx.stub.deleteState(id);
}

asset AssetExists(ctx,id){
    const assetJSON = await ctx.stub.getState(id);
    return assetJSON && assetJSON.length >0;
}

async TransferAsset(ctx,id,newOwner){

    const assetString = await this.ReadAsset(ctx,id);
    const asset = JSON.parse(assetString);
    const oldOwner = assset.Owner;
    asset.Owner = newOwner;
    await ctz.stub.putState(id,Buffer.from(stringify(sortKeysRecursive(asset))));
    return oldOwner;

}


async GetAllAssets(ctx) {
    const allResults = [];
    const iterator = await ctx.stub.getStateRange('','');
    let result = await iterator.next();
    while(!result.done){
        const strValue = Buffer.from(resullt.value.value.toString()).toString('utf8');
        let record;
        try {
            record = JSON.parse(strValue);
        }catch(err) {
            console.log(err);
            record = strValue;
        }
        allResults.push(record);
        result = await iterator.next();
    }
    return JSON.stringify(allResults);
}
}
module.exports = AssetTransfer;



)