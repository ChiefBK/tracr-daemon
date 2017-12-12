package command

import (
	log "github.com/inconshreveable/log15"
	"goku-bot/command/conditions"
	"goku-bot/command/actions"
	"io/ioutil"
	"encoding/json"
	"errors"
)

var bots []*Bot

func Init() {
	log.Info("Initializing command module", "module", "command")
	log.Debug("Creating bots", "module", "command")

	// initialize condition function map
	conditions.ConditionFunctions["TrueFunction"] = conditions.TrueFunction

	// initialize action function map
	actions.ActionFunctions["ShortPositionAction"] = actions.ShortPositionAction

	//path := filepath.Join("bot_templates", "bot1Template.json")
	//bot, err := readBotFile(path)

	//if err != nil {
	//	log.Error("there was an error reading the bot file", "module", "command", "file", path, "error", err)
	//}

	//log.Debug("Bot created", "module", "command", "bot", bot)
	//log.Debug("strategies created", "module", "command", "strats", bot.strategies)
	//log.Debug("trees created", "module", "command", "trees", bot.strategies["CLOSED"].decisionTrees)
	//log.Debug("root created", "module", "command", "signals", bot.strategies["CLOSED"].decisionTrees[0].root)
	//log.Debug("child created", "module", "command", "signals", bot.strategies["CLOSED"].decisionTrees[0].root.children[0])

	//rootSignal := NewSignal(conditions.TrueFunction, nil, true)
	//leafSignal := NewSignal(conditions.TrueFunction, actions.ShortPositionAction(), false)
	//tree := BuildDecisionChain(CLOSED_POSITION, rootSignal, leafSignal)
	//addBot("bot1", "poloniex", "USDT_BTC", nil, tree)

}

func Start() {
	log.Info("Starting command module", "module", "command")
	for _, bot := range bots {
		go bot.start()
	}
}

func readBotFile(filePath string) (*Bot, error) {
	rawJson, _ := ioutil.ReadFile(filePath)
	var data map[string]interface{}

	err := json.Unmarshal(rawJson, &data)

	if err != nil {
		//log.Error("there was an error un-marshalling stratagies file", "module", "command", "file", filePath)
		return nil, errors.New("there was an error un-marshalling strategies file")
	}

	name := data["name"].(string)
	pair := data["pair"].(string)
	exchange := data["exchange"].(string)
	bot := NewBot(name, exchange, pair)
	strategies := data["strategies"].([]interface{})

	for _, strategy := range strategies {
		position := strategy.(map[string]interface{})["position"].(string)
		trees := strategy.(map[string]interface{})["trees"].([]interface{})
		strat := NewStategy(position)

		for _, tree := range trees {
			rootSignal := tree.(map[string]interface{})["root"]
			signal := buildRoot(rootSignal.(map[string]interface{}))
			decisionTree := newDecisionTree(signal)
			strat.AddTree(decisionTree)
		}

		bot.addStrategy(strat)
	}

	return bot, nil
}

func buildRoot(root map[string]interface{}) *Signal {
	signal := createSignalFromInterface(root)
	children := root["children"].([]interface{})
	for _, child := range children {
		signal.addChild(buildRoot(child.(map[string]interface{})))
	}
	return signal
}

func createSignalFromInterface(raw map[string]interface{}) *Signal {
	isRoot := raw["isRoot"].(bool)
	conditionFunctionName := raw["condition"].(string)
	actionFunctionName, actionNotNull := raw["action"].(string)

	condFunc, ok := conditions.ConditionFunctions[conditionFunctionName]

	if !ok {
		// handle error
	}

	var action *actions.Action
	if actionNotNull {
		actionFunc, ok := actions.ActionFunctions[actionFunctionName]

		if !ok {
			// handle error
		}

		action = actionFunc()
	} else {
		action = nil
	}

	signal := NewSignal(condFunc, action, isRoot)
	return signal
}
