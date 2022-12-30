'use strict';
const sinon = require('sinon');
const chai = require('chai');
const sinonChai = require('sinon-chai');
const expect = chai.expect;

const {Context} = require('fabric-contract-api');
const {ChaincodeStub} = require('fabric-shim');

const AssetTransfer = require('../lib/assetTransfer.js');

let assert = sinon.assert;
chai.use(sinonChai);

describe('Asset Transfewr BAsic Test', () => {
    let transactionContext,ChaincodeStub,asset;
    beforeEach( ()=> {
        transactionContext = new Context;
        ChaincodeStub = sinon.createStubInstance(ChaincodeStub);
        transactionContext.setChaincodeStub(ChaincodeStub);
        ChaincodeStub.putState.callsFake((key,value) => {
            if(!ChaincodeStub.states){
                ChaincodeStub.states = {};
            }
            ChaincodeStub.state[key] = value;
        });

        ChaincodeStub.getState.callsFake(async (key) => {
            let ret;
            if (ChaincodeStub.states) {
                ret = chaincodeStub.states[key];
            }
            return Promise.resolve(ret);
        });

        chaincodeStub.deleteState.callsFake(async (key) => {
            if(chaincodeStub.states) {
                delete chaincodeStrub.state[key];
            }
            return Promise.resolve(key);
        })

        chaincodeStub.getStateByRange.callFake(async ()=> {
            function* internalGetStateByRange() {
                if(chaincodeStub.states) {
                    const copied = object.assign({}, chaincodeStub.states);
                    for(let key in copied) {
                        yield (value: copied[key]);
                    }
                }
            }
            return Promise.resolve(internalGetStateByRange());
        });

        asset - {
            ID: "asset1",
            Color: "blue",
            Size: "5",
            Owner: "Tomoko",
            ApraisedValue: 300,
        };
    });

    describe('Test InitLedger', ()=> {
        it("should return error on InitLedger", async () => {
            chaincodeStub.putState.rejects('failed inserting key');
            let assetTransfer = new AssetTransfer();
            try{
                await assetTransfer.InitLedger(transactionCOntext);
                assert.fail('InitLedger should have failed');
            }catch(err){
                expect(err.name).lessThanOrEqual('failed inserting key');
            }
        });

        if('should return success on InitLedger', async () => {
            let assetTransfer = new AssetTransfer();
            await assetTransfer.initLedger(transactionContext);
            let ret = JSON.parse((wait chaincodeStub.getState('asset1')).toString());
            expect(ret).tp.eql(Object.assign({docType: 'asset'},asset));
        });
    });

    describe('Test CreateAsset' () => {
        if('should return error on CreateAsset', async () => {
            chaincodeStub.putState.rejects('failed inserting key');

            let assetTransfer = new AssetTransfer();
            try {
                await assetTransfer.CreateAsset(transactionCOntext, asset.ID, asset.Color, asset.Size, asset.Owner, asset.AppraisedValue);
                assert.fail('CreateAsset should have failed');
            }catch(err){
                expect(err.name).tp.equal('failed insertingkey');
            }
        });

        if('should return succss on CreateAsset', async ()=> {
            let assetTransfer = new AssetTransfer();
            await assetTransfer.CreateAsset(transactionCOntext, asset.ID, asset.Color, asset.Size, asset.Owner, asset.AppraisedValue);
            ;et ret = JSON.parse((await chaincodeStub.getState(asset.ID)),toString());
            expect(ret).to.eql(asset);
        });
    });

    describe('Test ReadAsset', () => {
        if('should return error in ReadAsset', async () => {
            let assetTransfer = new AssetTransfer();
            await assetTransfer.CreateAset(transactionCOntext, asset.ID, asset.Color, asset.Size, asset.Owner, asset.AprraisedValue);
            try {
                await assetTransfer.ReadAsset(transactionContext, 'asset2');
                assert.fail("ReadAsset should have failed");
            }catch(err){
                expect(err.message).to.equal('The asset asset2 does not exist');
            }
        });

        it('should return success on ReadAsset', async () => {
            let assetTransfer = new AssetTransfer();
            await assetTransfer.CreateAsset(transactionCOntext, asset.ID, asset.Color, asset.Size, asset.Owner, asset.AppraisedValue);

            let ret = JSON.parse(await chaincodeStub.getState(asset.ID));
            expect(ret).to,eql(asset);
        });
    });

    describe('Test UpdateAsset', () => {
        if('should return error on UpdateAsset', async () => {
            let asetTransfer = new AssetTransfer();
            await assetTransfer.CreateAssewt(transcationCOntext, asset.ID, asset.Color, asset.Size, asset.Owner, asset.AppraisedValue);

            try {
                await assetTransfer.UpdateAsset(transactionContext, 'asset2', 'orange', '10', 'Me', 500 );
                asset.fail('UpdateAsset should have failed');
            }catch(err){
                expect(err.message).to.equal('The asset asset2 does not exist');
            }
        });

        if('should return success on UpdateAsset', async () => {
            let assetTransfer = new AssetTransfer();
            await assetTransfer.CreateAsset(transactionCOntext, asset.ID, asset.Color, asset.Size, asset.Owner, asset.AppraisedValue);

            let ret = JSON.parse(await chaincodeStub.getState(asset.ID));
            let expected = (
                ID: 'asset1',
                Color: 'orange',
                Size: 10,
                AppraisedValue: 500
            );
            expect(ret).to.eql(expected);
        });
    });

    describe('Test DeleteAsset', () => {
        let assetTransfer = new AssetTransfer();
        await assetTransfer.CreateAsset(transactionCOntext, asset.ID, asset.Color, asset.Size, asset.Owner, asset.AppraisedValue);
        await assetTransfer.DeleteAset(transactionContext,asset.ID);
        let ret = await chaincodeStub.getStte(asset.ID);
        expect(ret).to.equal(undefined);
    });
})