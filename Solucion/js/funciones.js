function drawSteps(pasos)
{
	var id_canvas;
	for(let i = 0; i < pasos.length; i++)
	{
		id_canvas = "canvas-gen-"+(i+1);
		drawPoints(pasos[i], id_canvas);
		drawPath(pasos[i], id_canvas);
	}
}

function drawBestSolution(puntos)
{
	drawPoints(puntos, "canvas-best-sol");
	drawPath(puntos, "canvas-best-sol")
}

function drawPoints(puntos, id_canvas)
{
	canvas = document.getElementById(id_canvas);
	ctx = canvas.getContext('2d');

	ctx.fillStyle = COLOR_FONDO;
	ctx.fillRect(0,0,canvas.width,canvas.height);

	ctx.fillStyle = COLOR_PUNTOS;
	puntos.forEach(punto => {
		ctx.beginPath();
		ctx.arc(punto.x, canvas.height - punto.y, 5, 0, 2 * Math.PI);
		ctx.fill();
	});
}

function drawPath(path, id_canvas)
{
	// Cargando imagenes
	img_start = new Image();
	img_finish = new Image();
	img_start.src = PATH_IMG_START;
	img_finish.src = PATH_IMG_FINISH;

	canvas = document.getElementById(id_canvas);
	ctx = canvas.getContext('2d');

	ctx.strokeStyle = COLOR_LINEAS;
	ctx.lineWidth = 2;
	var num_puntos = path.length;
	for (var i = 1; i < num_puntos; i++)
	{
		ctx.beginPath();
		ctx.moveTo(path[i-1].x, canvas.height - path[i-1].y);
		ctx.lineTo(path[i].x, canvas.height - path[i].y);
		ctx.stroke();
	}
	ctx.beginPath();
	ctx.moveTo(path[0].x, canvas.height - path[0].y);
	ctx.lineTo(path[num_puntos-1].x, canvas.height - path[num_puntos-1].y);
	ctx.stroke();

	img_start.onload = (e) => {
		canvas = document.getElementById(id_canvas);
		ctx = canvas.getContext('2d');
		ctx.drawImage(img_start, path[0].x - (img_start.width / 2), canvas.height - path[0].y - (img_start.height / 2));
	};
	img_finish.onload = (e) => {
		canvas = document.getElementById(id_canvas);
		ctx = canvas.getContext('2d');
		ctx.drawImage(img_finish, path[num_puntos - 1].x - 10, canvas.height - path[num_puntos - 1].y - img_finish.height);
	}
}

function drawAvanceOptimo()
{
	var data = new google.visualization.DataTable();

	data.addColumn('number', 'Iteraciones');
	data.addColumn('number', 'CostoDeRuta');

	data.addRows(mtz_avance_opt);

	var options = {
		height: 600,
		colors: ['#3bafda', 'blue', '#3fc26b'],
		hAxis: {title: 'Iteraciónes'},
		vAxis: {title: 'Costo de Ruta'}
	};

	var grafica = new google.visualization.AreaChart(document.getElementById('chart-avance-optimo'));
	grafica.draw(data, options);
}