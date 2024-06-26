package utils

import (
	"bytes"
	"context"
	"fmt"
	"io"

	// must be imported to init() raw-codec support
	dagpb "github.com/ipld/go-codec-dagpb"
	_ "github.com/ipld/go-ipld-prime/codec/raw"
	"github.com/ipld/go-ipld-prime/linking"
	"github.com/ipld/go-ipld-prime/node/basicnode"

	"github.com/ipfs/go-cid"
	mdagipld "github.com/ipfs/go-ipld-format"
	"github.com/ipfs/go-unixfsnode"
	"github.com/ipld/go-ipld-prime"
	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
	"github.com/ipld/go-ipld-prime/traversal"
	"github.com/ipld/go-ipld-prime/traversal/selector"
	selectorparse "github.com/ipld/go-ipld-prime/traversal/selector/parse"
)

func CreateLinkSystem(ds mdagipld.DAGService) linking.LinkSystem {
	// this is how we implement GETs
	linkSystem := cidlink.DefaultLinkSystem()
	linkSystem.StorageReadOpener = func(lctx ipld.LinkContext, lnk ipld.Link) (io.Reader, error) {
		cl, isCid := lnk.(cidlink.Link)
		if !isCid {
			return nil, fmt.Errorf("unexpected link type %#v", lnk)
		}

		node, err := ds.Get(lctx.Ctx, cl.Cid)
		if err != nil {
			return nil, err
		}

		return bytes.NewBuffer(node.RawData()), nil
	}
	unixfsnode.AddUnixFSReificationToLinkSystem(&linkSystem)
	return linkSystem
}

func TraverseDag(
	ctx context.Context,
	ds mdagipld.DAGService,
	startFrom cid.Cid,
	optionalSelector ipld.Node,
	visitCallback traversal.AdvVisitFn,
) error {

	if optionalSelector == nil {
		optionalSelector = selectorparse.CommonSelector_MatchAllRecursively
	}

	parsedSelector, err := selector.ParseSelector(optionalSelector)
	if err != nil {
		return err
	}

	linkContext := ipld.LinkContext{Ctx: ctx}
	linkSystem := CreateLinkSystem(ds)
	// this is what allows us to understand dagpb
	nodePrototypeChooser := dagpb.AddSupportToChooser(
		func(ipld.Link, ipld.LinkContext) (ipld.NodePrototype, error) {
			return basicnode.Prototype.Any, nil
		},
	)

	// this is how we pull the start node out of the DS
	startLink := cidlink.Link{Cid: startFrom}
	startNodePrototype, err := nodePrototypeChooser(startLink, linkContext)
	if err != nil {
		return err
	}
	startNode, err := linkSystem.Load(
		linkContext,
		startLink,
		startNodePrototype,
	)
	if err != nil {
		return err
	}

	// this is the actual execution, invoking the supplied callback
	return traversal.Progress{
		Cfg: &traversal.Config{
			Ctx:                            ctx,
			LinkSystem:                     linkSystem,
			LinkTargetNodePrototypeChooser: nodePrototypeChooser,
		},
	}.WalkAdv(startNode, parsedSelector, visitCallback)
}
